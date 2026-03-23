package usecase

import (
	"context"
	"duck-dns-changer/configs"
	"duck-dns-changer/internal/domain"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	ipifyEndpoint    = "https://api.ipify.org?format=json"
	duckDNSEndpoint  = "https://www.duckdns.org/update"
	maxResponseBytes = 8 * 1024
)

type IPHandler interface {
	GetIP(ctx context.Context) (domain.Ip, error)
	UpdateIP(ctx context.Context, ip domain.Ip) (bool, error)
}

type IPUseCase struct {
	cfg    configs.Config
	client *http.Client
}

func New(cfg configs.Config, client *http.Client) IPHandler {
	if client == nil {
		client = http.DefaultClient
	}

	return &IPUseCase{
		cfg:    cfg,
		client: client,
	}
}

func (i *IPUseCase) GetIP(ctx context.Context) (domain.Ip, error) {
	var ip domain.Ip

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ipifyEndpoint, nil)
	if err != nil {
		return ip, fmt.Errorf("create ipify request: %w", err)
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return ip, fmt.Errorf("call ipify endpoint: %w", err)
	}
	defer resp.Body.Close()

	bodyText, err := readResponseBody(resp.Body)
	if err != nil {
		return ip, fmt.Errorf("read ipify response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return ip, fmt.Errorf("ipify returned %s: %s", resp.Status, bodyText)
	}

	if err := json.Unmarshal([]byte(bodyText), &ip); err != nil {
		return ip, fmt.Errorf("decode ipify JSON: %w", err)
	}

	if strings.TrimSpace(ip.Ip) == "" {
		return ip, fmt.Errorf("ipify response has empty ip field")
	}

	return ip, nil
}

func (i *IPUseCase) UpdateIP(ctx context.Context, ip domain.Ip) (bool, error) {
	if i.cfg.Token == "" || i.cfg.Domain == "" || ip.Ip == "" {
		return false, fmt.Errorf("one of the required params is missing")
	}

	params := url.Values{}
	params.Set("domains", i.cfg.Domain)
	params.Set("token", i.cfg.Token)
	params.Set("ip", ip.Ip)
	params.Set("verbose", "true")

	endpoint := fmt.Sprintf("%s?%s", duckDNSEndpoint, params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return false, fmt.Errorf("create duckdns request: %w", err)
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("call duckdns endpoint: %w", err)
	}
	defer resp.Body.Close()

	responseText, err := readResponseBody(resp.Body)
	if err != nil {
		return false, fmt.Errorf("read duckdns response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("duckdns returned %s: %s", resp.Status, responseText)
	}

	if !isDuckDNSSuccess(responseText) {
		return false, fmt.Errorf("duckdns update failed: %s", responseText)
	}

	return true, nil
}

func readResponseBody(body io.Reader) (string, error) {
	b, err := io.ReadAll(io.LimitReader(body, maxResponseBytes))
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(b)), nil
}

func isDuckDNSSuccess(response string) bool {
	if response == "" {
		return false
	}

	for _, line := range strings.Split(response, "\n") {
		if strings.TrimSpace(line) == "OK" {
			return true
		}
	}

	return false
}
