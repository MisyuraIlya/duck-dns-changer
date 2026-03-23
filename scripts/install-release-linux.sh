#!/usr/bin/env bash
set -euo pipefail

service_name="duck-dns-changer"
install_dir="/usr/local/bin"
binary_path="${install_dir}/${service_name}"
systemd_dir="/etc/systemd/system"
etc_dir="/etc/${service_name}"
env_file="${etc_dir}/${service_name}.env"
run_time="03:00:00"
enable_timer=true
run_now=false
domain=""
token=""

usage() {
  cat <<'EOF'
Usage:
  sudo ./install.sh [options]

Options:
  --domain VALUE      DuckDNS domain (without .duckdns.org)
  --token VALUE       DuckDNS token
  --run-time HH:MM:SS Daily execution time (default: 03:00:00)
  --run-now           Execute service once immediately after install
  --no-enable         Install files but do not enable/start timer
  --help              Show help
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --domain)
      domain="${2:-}"
      shift 2
      ;;
    --token)
      token="${2:-}"
      shift 2
      ;;
    --run-time)
      run_time="${2:-}"
      shift 2
      ;;
    --run-now)
      run_now=true
      shift
      ;;
    --no-enable)
      enable_timer=false
      shift
      ;;
    --help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown argument: $1"
      usage
      exit 1
      ;;
  esac
done

if [[ "${EUID}" -ne 0 ]]; then
  echo "Run as root: sudo ./install.sh [options]"
  exit 1
fi

if ! command -v systemctl >/dev/null 2>&1; then
  echo "systemctl is not available on this machine."
  exit 1
fi

if [[ ! "${run_time}" =~ ^([01][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$ ]]; then
  echo "--run-time must be in HH:MM:SS format (24-hour clock)."
  exit 1
fi

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
src_binary="${script_dir}/${service_name}"
src_service="${script_dir}/${service_name}.service"
src_timer="${script_dir}/${service_name}.timer"

for file in "${src_binary}" "${src_service}" "${src_timer}"; do
  if [[ ! -f "${file}" ]]; then
    echo "Required file is missing next to install.sh: ${file}"
    exit 1
  fi
done

if [[ ! -x "${src_binary}" ]]; then
  chmod +x "${src_binary}"
fi

install -d -m 0755 "${install_dir}"
install -m 0755 "${src_binary}" "${binary_path}"

install -d -m 0750 "${etc_dir}"

existing_domain=""
existing_token=""
if [[ -f "${env_file}" ]]; then
  existing_domain="$(grep -E '^DOMAIN=' "${env_file}" | tail -n 1 | cut -d= -f2- || true)"
  existing_token="$(grep -E '^TOKEN=' "${env_file}" | tail -n 1 | cut -d= -f2- || true)"
fi

if [[ -z "${domain}" ]]; then
  domain="${existing_domain}"
fi

if [[ -z "${token}" ]]; then
  token="${existing_token}"
fi

if [[ -z "${domain}" ]]; then
  domain="your-duckdns-domain"
fi

if [[ -z "${token}" ]]; then
  token="your-duckdns-token"
fi

cat >"${env_file}" <<EOF
DOMAIN=${domain}
TOKEN=${token}
EOF
chmod 0640 "${env_file}"

install -m 0644 "${src_service}" "${systemd_dir}/${service_name}.service"

timer_tmp="$(mktemp)"
cp "${src_timer}" "${timer_tmp}"
sed -i "s/^OnCalendar=.*/OnCalendar=*-*-* ${run_time}/" "${timer_tmp}"
install -m 0644 "${timer_tmp}" "${systemd_dir}/${service_name}.timer"
rm -f "${timer_tmp}"

systemctl daemon-reload

if [[ "${enable_timer}" == "true" ]]; then
  systemctl enable --now "${service_name}.timer"
fi

if [[ "${run_now}" == "true" ]]; then
  systemctl start "${service_name}.service"
fi

echo "Installed ${service_name}."
echo "Binary: ${binary_path}"
echo "Environment file: ${env_file}"
echo "Schedule: daily at ${run_time}"

if [[ "${domain}" == "your-duckdns-domain" || "${token}" == "your-duckdns-token" ]]; then
  echo "WARNING: edit ${env_file} and set real DOMAIN/TOKEN."
fi

echo "Useful commands:"
echo "  systemctl status ${service_name}.timer"
echo "  systemctl list-timers ${service_name}.timer"
echo "  journalctl -u ${service_name}.service -n 50 --no-pager"
