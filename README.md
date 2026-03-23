# duck-dns-changer

## Installation (Ubuntu)

1. Download release package (`amd64` example):

```bash
curl -fL -o duck-dns-changer.tar.gz \
  https://github.com/MisyuraIlya/duck-dns-changer/releases/download/v1.0.0/duck-dns-changer_linux_amd64.tar.gz
```

For ARM64 server, use:

```bash
curl -fL -o duck-dns-changer.tar.gz \
  https://github.com/MisyuraIlya/duck-dns-changer/releases/download/v1.0.0/duck-dns-changer_linux_arm64.tar.gz
```

2. Extract and install:

```bash
tar -xzf duck-dns-changer.tar.gz
sudo ./install.sh --domain your-duckdns-domain --token your-duckdns-token --run-now
```

3. Verify timer:

```bash
sudo systemctl status duck-dns-changer.timer
sudo systemctl list-timers duck-dns-changer.timer
```

4. Check service logs:

```bash
sudo journalctl -u duck-dns-changer.service -n 50 --no-pager
```

Optional: choose custom daily run time:

```bash
sudo ./install.sh --domain your-duckdns-domain --token your-duckdns-token --run-time 04:30:00
```
