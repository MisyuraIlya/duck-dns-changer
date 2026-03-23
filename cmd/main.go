package main

import (
	"context"
	"duck-dns-changer/configs"
	"duck-dns-changer/internal/usecase"
	"log"
	"net/http"
	"os"
	"time"
)

const requestTimeout = 15 * time.Second

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.LUTC)

	cfg, err := configs.New()
	if err != nil {
		logger.Fatalf("load config: %v", err)
	}

	httpClient := &http.Client{
		Timeout: requestTimeout,
	}

	ipHandler := usecase.New(cfg, httpClient)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	ip, err := ipHandler.GetIP(ctx)
	if err != nil {
		logger.Fatalf("get public IP: %v", err)
	}

	updated, err := ipHandler.UpdateIP(ctx, ip)
	if err != nil {
		logger.Fatalf("update duckdns: %v", err)
	}

	if !updated {
		logger.Fatalf("duckdns update returned unsuccessful status")
	}

	logger.Printf("DuckDNS updated: domain=%s ip=%s", cfg.Domain, ip.Ip)
}
