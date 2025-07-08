#!/bin/bash

set -e

if ! [ -t 0 ]; then
  echo "This script must be run in an interactive shell. Please download and run it directly."
  exit 1
fi

DOCKER_COMPOSE_URL="https://raw.githubusercontent.com/panelium/panelium/main/assets/docker-compose.yml"
SYSTEMD_BASE_URL="https://raw.githubusercontent.com/panelium/panelium/main/assets"
NGINX_CONF_URL="https://raw.githubusercontent.com/panelium/panelium/main/assets/panelium.conf"

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

# Check for nginx
if ! command -v nginx &>/dev/null; then
  echo "Nginx is not installed. Please install Nginx first." >&2
  exit 1
fi

# Check for certbot
if ! command -v certbot &>/dev/null; then
  echo "Certbot is not installed. Please install Certbot first." >&2
  exit 1
fi

# Check for jq
if ! command -v jq &>/dev/null; then
  echo "jq is not installed. Please install jq first." >&2
  exit 1
fi

# Check for existing install unless --unsupported-ignore-existing is passed
IGNORE_EXISTING=0
for arg in "$@"; do
  if [[ "$arg" == "--unsupported-ignore-existing" ]]; then
    IGNORE_EXISTING=1
  fi
done

if [[ -d /etc/panelium/backend || -d /etc/panelium/daemon || -d /etc/panelium/dashboard || -f /var/lib/panelium/docker-compose.yml ]]; then
  if [[ $IGNORE_EXISTING -eq 0 ]]; then
    echo "Panelium appears to be already installed."
    echo "To override, run with --unsupported-ignore-existing. This flag is NOT SUPPORTED and may BREAK your install. Use at your own risk!"
    exit 1
  else
    echo "Warning: Existing Panelium install detected. Proceeding due to --unsupported-ignore-existing flag."
    echo "This flag is NOT SUPPORTED and may BREAK your install. Use at your own risk!"
  fi
fi

# Prompt for domains (used for both URLs and certbot/nginx)
read -rp "Enter the Dashboard domain (e.g. dashboard.example.com): " DASHBOARD_DOMAIN
read -rp "Enter the Backend domain (e.g. backend.example.com): " BACKEND_DOMAIN
read -rp "Enter the Daemon domain (e.g. daemon.example.com): " DAEMON_DOMAIN

DASHBOARD_URL="https://$DASHBOARD_DOMAIN"
BACKEND_URL="https://$BACKEND_DOMAIN"
DAEMON_URL="https://$DAEMON_DOMAIN"

# Create config directories if not exist
mkdir -p /etc/panelium/backend
mkdir -p /etc/panelium/daemon
mkdir -p /etc/panelium/dashboard
mkdir -p /var/lib/panelium

# Download docker-compose.yml
if ! curl -fsSL "$DOCKER_COMPOSE_URL" -o /var/lib/panelium/docker-compose.yml; then
  echo "Failed to download docker-compose.yml from $DOCKER_COMPOSE_URL" >&2
  exit 1
fi

# Update backend config
BACKEND_CONFIG="/etc/panelium/backend/config.json"
if [[ -f $BACKEND_CONFIG ]]; then
  jq \
    --arg dashboard "$DASHBOARD_URL" \
    --arg backend "$BACKEND_URL" \
    '.hosts.dashboard = $dashboard | .hosts.backend = $backend' \
    "$BACKEND_CONFIG" > "$BACKEND_CONFIG.tmp" && mv "$BACKEND_CONFIG.tmp" "$BACKEND_CONFIG"
else
  cat > "$BACKEND_CONFIG" <<EOF
{
  "hosts": {
    "dashboard": "$DASHBOARD_URL",
    "backend": "$BACKEND_URL"
  }
}
EOF
fi

# Update daemon config
DAEMON_CONFIG="/etc/panelium/daemon/config.json"
if [[ -f $DAEMON_CONFIG ]]; then
  jq \
    --arg dashboard "$DASHBOARD_URL" \
    --arg backend "$BACKEND_URL" \
    --arg daemon "$DAEMON_URL" \
    '.hosts.dashboard = $dashboard | .hosts.backend = $backend | .hosts.daemon = $daemon' \
    "$DAEMON_CONFIG" > "$DAEMON_CONFIG.tmp" && mv "$DAEMON_CONFIG.tmp" "$DAEMON_CONFIG"
else
  cat > "$DAEMON_CONFIG" <<EOF
{
  "hosts": {
    "dashboard": "$DASHBOARD_URL",
    "backend": "$BACKEND_URL",
    "daemon": "$DAEMON_URL"
  }
}
EOF
fi

# Update dashboard config
DASHBOARD_CONFIG="/etc/panelium/dashboard/config.json"
cat > "$DASHBOARD_CONFIG" <<EOF
{
  "BACKEND_HOST": "$BACKEND_URL",
  "TURNSTILE_SITE_KEY": "1x00000000000000000000AA"
}
EOF

# Obtain SSL certificates using certbot (standalone mode)
certbot certonly --nginx --non-interactive --agree-tos --register-unsafely-without-email -d "$DASHBOARD_DOMAIN" -d "$BACKEND_DOMAIN" -d "$DAEMON_DOMAIN"

# Download and configure nginx site config
NGINX_CONF_PATH="/etc/nginx/sites-available/panelium.conf"
if ! curl -fsSL "$NGINX_CONF_URL" -o "$NGINX_CONF_PATH"; then
  echo "Failed to download nginx config from $NGINX_CONF_URL" >&2
  exit 1
fi

# Replace example.com domains in nginx config with user input

get_cert_folder() {
  # Extract the third-level domain for cert folder if more levels are present
  domain="$1"
  IFS='.' read -ra parts <<< "$domain"
  count=${#parts[@]}
  if (( count > 3 )); then
    # Use the last three parts (e.g. backend.panelium-demo.ndmh.xyz -> panelium-demo.ndmh.xyz)
    echo "${parts[count-3]}.${parts[count-2]}.${parts[count-1]}"
  else
    echo "$domain"
  fi
}

escape_sed() {
  # Escape &, /, and newlines for sed replacement
  printf '%s' "$1" | sed -e 's/[&/\\]/\\&/g' | tr '\n' '\\n'
}

dashboard_cert_folder=$(get_cert_folder "$DASHBOARD_DOMAIN")
backend_cert_folder=$(get_cert_folder "$BACKEND_DOMAIN")
daemon_cert_folder=$(get_cert_folder "$DAEMON_DOMAIN")

# Update certificate paths in nginx config to use the correct cert folder
sed -i "s|/etc/letsencrypt/live/dashboard.example.com/|/etc/letsencrypt/live/$dashboard_cert_folder/|g" "$NGINX_CONF_PATH"
sed -i "s|/etc/letsencrypt/live/backend.example.com/|/etc/letsencrypt/live/$backend_cert_folder/|g" "$NGINX_CONF_PATH"
sed -i "s|/etc/letsencrypt/live/daemon.example.com/|/etc/letsencrypt/live/$daemon_cert_folder/|g" "$NGINX_CONF_PATH"

dashboard_escaped=$(escape_sed "$DASHBOARD_DOMAIN")
backend_escaped=$(escape_sed "$BACKEND_DOMAIN")
daemon_escaped=$(escape_sed "$DAEMON_DOMAIN")
sed -i "s/dashboard.example.com/$dashboard_escaped/g" "$NGINX_CONF_PATH"
sed -i "s/backend.example.com/$backend_escaped/g" "$NGINX_CONF_PATH"
sed -i "s/daemon.example.com/$daemon_escaped/g" "$NGINX_CONF_PATH"

# Enable nginx site
if [ -d /etc/nginx/sites-enabled ]; then
  ln -sf "$NGINX_CONF_PATH" /etc/nginx/sites-enabled/panelium.conf
else
  echo "Warning: /etc/nginx/sites-enabled does not exist. Please ensure your nginx includes /etc/nginx/sites-available/panelium.conf manually."
fi

# Test and reload nginx
nginx -t && systemctl reload nginx

# Detect OS
OS_TYPE=$(uname)

if [[ "$OS_TYPE" == "Darwin" ]]; then
  # macOS: print sample docker compose commands and create helper script
  echo "macOS detected. To manage Panelium services manually, use the following commands:"
  echo "Start all:   docker compose -f /var/lib/panelium/docker-compose.yml up -d"
  echo "Stop all:    docker compose -f /var/lib/panelium/docker-compose.yml down"
  echo "Start one:   docker compose -f /var/lib/panelium/docker-compose.yml up -d <service>"
  echo "Stop one:    docker compose -f /var/lib/panelium/docker-compose.yml stop <service>"

  cat > /var/lib/panelium/panelium-docker.sh <<EOS
#!/bin/sh
cd /var/lib/panelium
case "$1" in
  start)
    docker compose up -d ;;
  stop)
    docker compose down ;;
  restart)
    docker compose down && docker compose up -d ;;
  status)
    docker compose ps ;;
  *)
    echo "Usage: $0 {start|stop|restart|status}" ;;
esac
EOS
  chmod +x /var/lib/panelium/panelium-docker.sh
  echo "Helper script created at /var/lib/panelium/panelium-docker.sh."
  echo "Panelium configuration and setup complete!"
else
  # Linux: download and set up systemd unit files
  for unit in panelium.service paneliumb.service paneliumd.service; do
    if ! curl -fsSL "$SYSTEMD_BASE_URL/$unit" -o "/etc/systemd/system/$unit"; then
      echo "Failed to download $unit from $SYSTEMD_BASE_URL/$unit" >&2
      exit 1
    fi
  done

  echo "Reloading systemd daemon and enabling Panelium services..."
  systemctl daemon-reload
  systemctl enable panelium.service paneliumb.service paneliumd.service
  systemctl restart panelium.service paneliumb.service paneliumd.service
  echo "Panelium configuration and setup complete!"
fi

# Copy backend JWT public key to daemon config dir
BACKEND_JWT_PUBLIC_KEY="/etc/panelium/backend/jwt_public_key.pem"
DAEMON_JWT_PUBLIC_KEY="/etc/panelium/daemon/backend_jwt_public_key.pem"

# Wait up to 10 seconds for the backend JWT public key to appear
WAIT_TIME=0
while [[ ! -f "$BACKEND_JWT_PUBLIC_KEY" && $WAIT_TIME -lt 10 ]]; do
  echo "Waiting for backend JWT public key to be generated... ($(expr 10 - $WAIT_TIME)s remaining until exit)"
  sleep 1
  WAIT_TIME=$((WAIT_TIME+1))
done

if [[ -f "$BACKEND_JWT_PUBLIC_KEY" ]]; then
  cp "$BACKEND_JWT_PUBLIC_KEY" "$DAEMON_JWT_PUBLIC_KEY"
  echo "Copied backend JWT public key to daemon config directory."
  # Restart daemon service or container
  OS_TYPE=$(uname)
  if [[ "$OS_TYPE" == "Darwin" ]]; then
    echo "Restarting daemon container via docker compose..."
    docker compose -f /var/lib/panelium/docker-compose.yml restart daemon
  else
    echo "Restarting paneliumd.service via systemctl..."
    systemctl restart paneliumd.service
  fi
else
  echo "Warning: Backend JWT public key not found at $BACKEND_JWT_PUBLIC_KEY. Manual copy is required, otherwise the daemon will not start."
fi
