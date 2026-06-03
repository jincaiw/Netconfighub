FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o netconfighub ./cmd/api

FROM node:22-alpine AS frontend
WORKDIR /app/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ .
RUN npm run build

FROM alpine:3.19
RUN apk add --no-cache ca-certificates git git-crypt tzdata
WORKDIR /app
COPY --from=builder /app/netconfighub .
COPY --from=builder /app/configs ./configs
COPY --from=frontend /app/web/dist ./web/dist
RUN mkdir -p /app/data
EXPOSE 8080
ENTRYPOINT ["./netconfighub"]
