# This software has not yet been extensively tested for security and is therefore not recommended for production use.

# Panelium

Panelium is an all-in-one server hosting solution that provides a web-based control panel for managing games, servers,
applications, and services.

![Demo Screenshot](https://raw.githubusercontent.com/Panelium/Panelium/7556f9dcc11fa9e80b5dbd065d6c8363fffb7832/assets/panelium-screenshot.png)

## Installation

### Requirements

- Unix-like operating system (Linux, macOS, etc.), Windows is currently not supported.
- Docker
- Docker Compose
- Nginx
- Certbot
- jq

### Quick Start

You can install Panelium using the automated setup script. This script will:

- Check for required dependencies
- Download and configure Docker Compose and systemd units (Linux) or helper scripts (macOS)
- Set up Nginx and SSL certificates with Certbot
- Prompt you for your dashboard, backend, and daemon domains

**To run the setup script:**

```bash
bash <(curl -s https://raw.githubusercontent.com/panelium/panelium/main/assets/panelium-setup.sh)
```

Follow the prompts to enter your domain names (e.g., dashboard.example.com, backend.example.com, daemon.example.com).
The script will handle the rest.

### Configuration

Configuration files are created automatically in `/etc/panelium/` for each service. You can adjust these files as
needed after installation.

Note: backend and daemon(s) have to be on the same second-level domain (third- and lower level domains, aka subdomains,
can be different) for the CORS and cookies to work properly.

Blueprints are available at https://blueprints.ndmh.xyz/ ([source](https://github.com/Panelium/Blueprints))

### Useful Commands

- **Make a user an admin**:
  ```sh
  docker compose -f /var/lib/panelium/docker-compose.yml exec backend /apps/backend --make-admin <username/email>
  ```

### Uninstalling Panelium

To completely remove Panelium, including all configuration, data, and nginx configs, run the uninstall script:

```sh
bash <(curl -s https://raw.githubusercontent.com/panelium/panelium/main/assets/panelium-uninstall.sh)
```

You will be prompted for confirmation before any files are deleted.

### Development Environment Setup

To set up a development environment for Panelium, follow these steps:

#### 1. Install Dependencies

- Go (v1.24+)
- Node.js (v22+) and npm
- Docker & Docker Compose
- buf (for protobuf)

On Arch Linux:

- Note: You will need the AUR helper `yay` installed for this.

```sh
yay -Syu
yay -S go nodejs npm docker docker-compose buf
```

On Debian and Ubuntu or other derivatives:

```sh
sudo apt update
sudo apt install golang nodejs npm docker.io docker-compose buf
```

On Fedora:

```sh
sudo dnf install golang nodejs npm docker docker-compose buf
```

On macOS (with Homebrew):

```sh
brew install go node docker bufbuild/buf/buf
brew install --cask docker # for Docker Desktop
```

On Alpine:

```sh
apk add go nodejs npm docker docker-compose buf
```

> For Windows, WSL2 is recommended. We don't currently support native Windows installations.

#### 2. Install Go Tools

```sh
go install github.com/sudorandom/protoc-gen-connect-openapi@main
```

> **Note:** Make sure your Go bin directory (usually `$HOME/go/bin`) is in your `PATH` so installed Go tools are
> available.

#### 3. Clone the Repository

```sh
git clone https://github.com/Panelium/Panelium.git
cd Panelium
```

#### 4. Build and Run with Docker Compose

For development, you can use the provided `docker-compose.yml` in the root directory:

```sh
docker compose up --build
```

#### 5. Manual Development Commands (via Makefile)

You can use the provided Makefile for common development tasks:

- Start the dashboard in dev mode (with live reload):
  ```sh
  make dev-dashboard
  ```
- Start the backend in dev mode (with live reload):
  ```sh
  make dev-backend
  ```
- Start the daemon in dev mode (with live reload):
  ```sh
  make dev-daemon
  ```
- Regenerate protobuf code:
  ```sh
  make gen-proto
  ```

#### 6. Protobuf/gRPC

If you make changes to protobuf files, regenerate code with:

```sh
cd shared/proto
buf generate
```

#### 7. Environment Variables & Config

- Configuration files are in `/etc/panelium/` for each service.
- For local development, you may want to adjust ports and hosts in the configs and docker-compose files.

#### 8. Useful Commands

- Rebuild containers: `docker compose build`
- Stop containers: `docker compose down`
- View logs: `docker compose logs -f`