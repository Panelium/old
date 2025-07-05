# Not ready for use yet, actively under development

Blueprints are available at http://blueprints.ndmh.xyz/ ([source](https://github.com/Panelium/Blueprints))

### Development Environment Setup

- Swap out the yay commands to equivalent commands of whatever package manager you use.
- Make sure you have your GOBIN in your PATH.

```
yay -S buf
go install github.com/sudorandom/protoc-gen-connect-openapi@main
```

Note: backend and daemon(s) have to be on the same domain (subdomains can be different) for the CORS and cookies to work
properly.