dev-dashboard:
	npm run dev -w dashboard

dev-backend:
	cd ./apps/backend && air

dev-daemon:
	cd ./apps/daemon && air

dashboard-image:
	docker build --target dashboard --tag nginxalpine_dashboard .

dashboard-container:
	docker run -p 80:80 nginxalpine_dashboard:latest

docker-dashboard: dashboard-image dashboard-container

gen-proto:
	cd ./shared/proto && buf generate

.PHONY: dev-dashboard dev-backend dev-daemon gen-proto
