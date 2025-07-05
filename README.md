# This software has not yet been extensively tested for security and is therefore not recommended for production use.

# Panelium

Panelium is an all-in-one server hosting solution that provides a web-based control panel for managing games, servers,
applications, and services.

## Installation

### Requirements

- Unix-like operating system (Linux, macOS, etc.), Windows is currently not supported.
- Docker
- Docker Compose

### Quick Start

To install Panelium, you can use the provided Docker Compose file. This will set up the necessary services, including
the frontend (dashboard), backend, and a daemon.

```bash
git clone https://github.com/Panelium/Panelium.git
cd Panelium
docker compose up -d
```

### Configuration

TBA

Note: backend and daemon(s) have to be on the same domain (subdomains can be different) for the CORS and cookies to work
properly.

Blueprints are available at http://blueprints.ndmh.xyz/ ([source](https://github.com/Panelium/Blueprints))

### Development Environment Setup
- Swap out the yay commands to equivalent commands of whatever package manager you use.
- Make sure you have your GOBIN in your PATH.

```
yay -S buf
go install github.com/sudorandom/protoc-gen-connect-openapi@main
```