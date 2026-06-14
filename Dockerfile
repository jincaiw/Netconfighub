# syntax=docker/dockerfile:1.7

FROM node:22-alpine AS frontend
WORKDIR /src/web
COPY web/package.json web/package-lock.json ./
RUN --mount=type=cache,target=/root/.npm npm ci
COPY web/ ./
RUN npm run build

FROM golang:1.25-alpine AS backend
ARG VERSION=dev
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
COPY --from=frontend /src/web/dist ./web/dist
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build \
    -buildvcs=false \
    -trimpath \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o /out/netconfighub ./cmd/api

FROM alpine:3.22
ARG VERSION=dev
LABEL org.opencontainers.image.title="NetConfigHub" \
      org.opencontainers.image.description="Network device configuration backup and management" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.source="https://github.com/jincaiw/Netconfighub"
RUN apk add --no-cache ca-certificates git git-crypt tzdata \
    && addgroup -S netconfighub \
    && adduser -S -G netconfighub -h /app netconfighub \
    && install -d -o netconfighub -g netconfighub /app/data
WORKDIR /app
COPY --from=backend --chown=netconfighub:netconfighub /out/netconfighub /usr/local/bin/netconfighub
USER netconfighub
ENV NCH_DATA_DIR=/app/data \
    NCH_SERVER_HOST=0.0.0.0 \
    NCH_SERVER_PORT=8080
VOLUME ["/app/data"]
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://127.0.0.1:8080/api/v1/health >/dev/null || exit 1
ENTRYPOINT ["/usr/local/bin/netconfighub"]
