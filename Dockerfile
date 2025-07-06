FROM golang:1.24-alpine AS build-backend
WORKDIR /app/backend/
RUN apk add --no-cache build-base
COPY apps/backend/ .
COPY shared/common/go/ /app/common/
COPY shared/protogen/go/ /app/proto_gen_go/
RUN go mod edit -replace=panelium/common=/app/common
RUN go mod edit -replace=panelium/proto_gen_go=/app/proto_gen_go
RUN go mod tidy
RUN go mod download
RUN go build -o /build/backend

FROM alpine:latest AS backend
EXPOSE 9090/tcp
WORKDIR /app/
COPY --from=build-backend /build/backend /app/backend
VOLUME ["/etc/panelium", "/var/lib/panelium", "/var/log/panelium"]
CMD ["./backend"]

FROM golang:1.24-alpine AS build-daemon
WORKDIR /app/daemon/
RUN apk add --no-cache build-base
COPY apps/daemon/ .
COPY shared/common/go/ /app/common/
COPY shared/protogen/go/ /app/proto_gen_go/
RUN go mod edit -replace=panelium/common=/app/common
RUN go mod edit -replace=panelium/proto_gen_go=/app/proto_gen_go
RUN go mod tidy
RUN go mod download
RUN go build -o /build/daemon

FROM alpine:latest AS daemon
EXPOSE 9000/tcp
WORKDIR /app/
COPY --from=build-daemon /build/daemon /app/daemon
VOLUME ["/etc/panelium", "/var/lib/panelium", "/var/log/panelium", "/var/run/docker.sock", "/var/lib/docker/volumes"]
CMD ["./daemon"]

FROM node:22-alpine AS build-dashboard
WORKDIR /app/dashboard/
COPY apps/dashboard/ .
RUN npm install
COPY shared/protogen/ts/ /app/node_modules/proto-gen-ts/
RUN npm install --prefix /app/node_modules/proto-gen-ts/
RUN npm run build

FROM nginx:alpine AS dashboard
EXPOSE 80/tcp
COPY --from=build-dashboard /app/dashboard/build/client/ /usr/share/nginx/html/
COPY assets/nginx.conf /etc/nginx/conf.d/default.conf
