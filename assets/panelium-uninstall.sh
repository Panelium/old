#!/bin/bash

set -e

if [[ $EUID -ne 0 ]]; then
  echo "This script must be run as root" >&2
  exit 1
fi

read -p "Are you sure you want to uninstall Panelium and remove all related files? (y/N): " confirm
if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
  echo "Uninstall cancelled."
  exit 0
fi

# Stop and disable systemd services if they exist
for svc in panelium.service paneliumb.service paneliumd.service; do
  if systemctl list-unit-files | grep -q "$svc"; then
    systemctl stop "$svc" || true
    systemctl disable "$svc" || true
    rm -f "/etc/systemd/system/$svc"
  fi
done
systemctl daemon-reload

# Stop docker compose if running (fallback)
if [ -f /var/lib/panelium/docker-compose.yml ]; then
  docker compose -f /var/lib/panelium/docker-compose.yml down || true
fi

# Remove config and data directories
rm -rf /etc/panelium /var/lib/panelium

# Remove nginx config
rm -f /etc/nginx/sites-available/panelium.conf
rm -f /etc/nginx/sites-enabled/panelium.conf

# Reload nginx
nginx -t && systemctl reload nginx || true

echo "Panelium has been uninstalled."

