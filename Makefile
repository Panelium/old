dev-dashboard:
	cd ./apps/dashboard && deno task dev

dev-backend:
	cd ./apps/backend && air

dev-daemon:
	cd ./apps/daemon && air

gen-proto:
	cd ./shared/proto && buf generate

.PHONY: dev-dashboard dev-backend dev-daemon gen-proto
