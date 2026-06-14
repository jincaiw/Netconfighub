#!/usr/bin/env sh
set -eu

VERSION="${VERSION:-v0.1.2}"
ROOT="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
DIST="$ROOT/dist"
STAGE="$DIST/netconfighub-$VERSION-linux-amd64"

rm -rf "$STAGE"
mkdir -p "$STAGE"

cd "$ROOT/web"
npm ci
npm run build

cd "$ROOT"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -buildvcs=false \
  -trimpath \
  -ldflags="-s -w -X main.version=$VERSION" \
  -o "$STAGE/netconfighub" ./cmd/api

cp README.md README.zh-CN.md LICENSE "$STAGE/"
cp deploy/netconfighub.service "$STAGE/"

tar -C "$DIST" -czf "$DIST/netconfighub-$VERSION-linux-amd64.tar.gz" \
  "netconfighub-$VERSION-linux-amd64"

(
  cd "$DIST"
  shasum -a 256 "netconfighub-$VERSION-linux-amd64.tar.gz" \
    > "netconfighub-$VERSION-linux-amd64.tar.gz.sha256"
)

file "$STAGE/netconfighub"
go version -m "$STAGE/netconfighub"
if [ "$(uname -s)" = "Linux" ]; then
  "$STAGE/netconfighub" -version
fi
cat "$DIST/netconfighub-$VERSION-linux-amd64.tar.gz.sha256"
