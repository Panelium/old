FROM golang:1.24-alpine AS build-backend
WORKDIR /app/backend/
RUN apk add --no-cache build-base
COPY apps/backend/ .
COPY shared/common/go/ /app/common/
COPY shared/protogen/go/ /app/proto-gen-go/
RUN go mod edit -replace=panelium/common=/app/common
RUN go mod edit -replace=panelium/proto-gen-go=/app/proto-gen-go
RUN go mod tidy
RUN go mod download
RUN go build -o /build/backend

FROM alpine:latest AS backend
ENV PORT=9090
EXPOSE ${PORT}/tcp
WORKDIR /app/
COPY --from=build-backend /build/backend /app/backend
CMD ["./backend"]
