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
curl -fsSL https://raw.githubusercontent.com/panelium/panelium/main/assets/panelium-setup.sh | bash
```

Follow the prompts to enter your domain names (e.g., dashboard.example.com, backend.example.com, daemon.example.com).
The script will handle the rest.

### Configuration

Configuration files are generated automatically in `/etc/panelium/` for each service. You can adjust these files as
needed after installation.

Note: backend and daemon(s) have to be on the same second-level domain (third- and lower level domains, aka subdomains,
can be different) for the CORS and cookies to work properly.

Blueprints are available at http://blueprints.ndmh.xyz/ ([source](https://github.com/Panelium/Blueprints))

### Development Environment Setup
- Swap out the yay commands to equivalent commands of whatever package manager you use. We have provided a commented out example for the apt package manager.
- Make sure you have your GOBIN in your PATH.

```
yay -Syu # apt update
yay -S buf # apt install buf
go install github.com/sudorandom/protoc-gen-connect-openapi@main
```
