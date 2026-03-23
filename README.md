# duck-dns-changer

Small Go app that updates your DuckDNS record with your current public IP.

## Build locally

```bash
go build -o duck-dns-changer ./cmd
```

## Create a GitHub release automatically

This repo includes a workflow: `.github/workflows/release.yml`.

It runs when you push a tag that starts with `v` (example: `v1.0.0`), builds binaries for:

- linux amd64
- linux arm64
- darwin amd64
- darwin arm64
- windows amd64

and uploads them to the GitHub Release page.

Release commands:

```bash
git tag v1.0.0
git push origin v1.0.0
```

## Run once per day on Ubuntu (`systemd`)

1. Download the Linux release package (replace owner/repo/version/arch):

```bash
curl -fL -o duck-dns-changer.tar.gz \
  https://github.com/<owner>/<repo>/releases/download/v1.0.0/duck-dns-changer_linux_amd64.tar.gz
```

2. Extract and install:

```bash
tar -xzf duck-dns-changer.tar.gz
sudo ./install.sh --domain your-duckdns-domain --token your-duckdns-token --run-now
```

This installs the binary, env file, service, and daily timer.

3. Check next scheduled execution:

```bash
sudo systemctl list-timers duck-dns-changer.timer
```

Default schedule is `03:00` daily. You can change it during install:

```bash
sudo ./install.sh --domain your-duckdns-domain --token your-duckdns-token --run-time 04:30:00
```
