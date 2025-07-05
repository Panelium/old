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
ENV PORT=9090
EXPOSE ${PORT}/tcp
WORKDIR /app/
COPY --from=build-backend /build/backend /app/backend
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
ENV PORT=9000
EXPOSE ${PORT}/tcp
WORKDIR /app/
COPY --from=build-daemon /build/daemon /app/daemon
VOLUME ["/var/run/docker.sock", "/var/lib/docker/volumes"]
# need to mount these when running the container
CMD ["./daemon"]

FROM node:22-alpine AS build-dashboard
WORKDIR /app/dashboard/
COPY apps/dashboard/ .
RUN npm install
# awful hack to not have to change imports
COPY shared/protogen/ts/ /app/node_modules/proto-gen-ts/
RUN npm install --prefix /app/node_modules/proto-gen-ts/
RUN npm run build

FROM nginx:alpine AS dashboard
COPY --from=build-dashboard /app/dashboard/build/ /usr/share/nginx/html/
