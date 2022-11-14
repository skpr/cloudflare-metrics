#!/usr/bin/make -f

export CGO_ENABLED=0

PROJECT=github.com/skpr/cloudflare-metrics
VERSION=$(shell git describe --tags --always)
COMMIT=$(shell git rev-list -1 HEAD)
IMAGE=skpr/cloudflare-metrics
VERSION=$(shell git describe --tags --always)

# Builds the project.
define go_build
	GOOS=${2} GOARCH=${3} go build -o bin/${1}_${2}_${3} -ldflags='-extldflags "-static" -X main.GitVersion=${VERSION}' ${4}
endef

# Builds the CLI.
build:
	$(call go_build,cloudflare_metrics,linux,amd64,${PROJECT})
	$(call go_build,cloudflare_metrics,darwin,amd64,${PROJECT})

# Run all lint checking with exit codes for CI.
lint:
	golint -set_exit_status `go list ./... | grep -v /vendor/`

# Run tests with coverage reporting.
test:
	go test -cover ./...


# Releases the project Docker Hub
release:
	docker build -t ${IMAGE}:${VERSION} -t ${IMAGE}:latest .
	docker push ${IMAGE}:${VERSION}
	docker push ${IMAGE}:latest

.PHONY: *
