VERSION ?= v0.1.2

.PHONY: test build-web build-linux release-bundle docker-build verify

test:
	go test ./...
	go vet ./...

build-web:
	cd web && npm ci && npm run build

build-linux: build-web
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-buildvcs=false -trimpath -ldflags="-s -w -X main.version=$(VERSION)" \
		-o dist/netconfighub-$(VERSION)-linux-amd64 ./cmd/api

release-bundle:
	VERSION=$(VERSION) ./scripts/package-release.sh

docker-build:
	docker build --build-arg VERSION=$(VERSION) \
		-t jincaiw/netconfighub:$(VERSION:v%=%) \
		-t jincaiw/netconfighub:latest .

verify: test build-web
	cd web && npx playwright test
