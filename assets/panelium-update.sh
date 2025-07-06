#!/bin/bash

set -e

if ! [ -t 0 ]; then
  echo "This script must be run in an interactive shell. Please download and run it directly."
  exit 1
fi

# Require root
if [[ $EUID -ne 0 ]]; then
  echo "This script must be run as root" >&2
  exit 1
fi

# Check for docker
if ! command -v docker &>/dev/null; then
  echo "Docker is not installed. Please install Docker first." >&2
  exit 1
fi

# Check for docker compose (plugin or standalone)
if ! docker compose version &>/dev/null && ! command -v docker-compose &>/dev/null; then
  echo "Docker Compose is not installed. Please install Docker Compose (plugin or standalone) first." >&2
  exit 1
fi

# Check if Panelium is installed
if [[ ! -f /var/lib/panelium/docker-compose.yml ]]; then
  echo "Panelium does not appear to be installed. Aborting update." >&2
  exit 1
fi

# Pull latest docker images
echo "Pulling latest docker images..."
docker compose -f /var/lib/panelium/docker-compose.yml pull

echo "Ready to restart Panelium services to apply updates."
read -rp "Proceed with restart? [y/N]: " CONFIRM
if [[ ! "$CONFIRM" =~ ^[Yy]$ ]]; then
  echo "Update pulled, but restart cancelled by user. Please restart services manually to apply updates."
  exit 0
fi

# Restart services
if [[ -f /etc/systemd/system/panelium.service ]]; then
  echo "Restarting Panelium systemd services..."
  systemctl restart panelium.service paneliumb.service paneliumd.service
else
  echo "Restarting Panelium docker containers..."
  docker compose -f /var/lib/panelium/docker-compose.yml up -d
fi

echo "Panelium update complete!"
