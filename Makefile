# Cargar variables de entorno desde .env si existe
ENV_LOAD = [ -f .env ] && set -a && . ./.env && set +a || true


ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

# Dependencies
install:
	go get ./...

install-playground:
	cd playground && \
	  bun install

# Tests
test:
	go test -cover -v ./...

test-json:
	go test -cover -json -v ./...

coverage:
	go test -coverprofile=tmp/coverage.out ./...
	go tool cover -html=tmp/coverage.out

# Linting
lint:
	test -z $(gofmt -l .)


# Development
api:
	$(ENV_LOAD); $(GOPATH)/bin/air -c .air-api.toml

run:
	$(ENV_LOAD); go run cmd/api/main.go

cli:
	$(ENV_LOAD); go run cmd/cli/main.go $(filter-out $@,$(MAKECMDGOALS))

dev:
	$(ENV_LOAD); go run ./cmd/dev/main.go $(filter-out $@,$(MAKECMDGOALS))

playground:
	$(ENV_LOAD); cd playground && \
	  bun dev


# Deployment
build-cli:
	CGO_ENABLED=0 go build -o bin/cli cmd/cli/main.go
	chmod +x bin/cli

build-api:
	CGO_ENABLED=0 go build -o bin/api cmd/api/main.go
	chmod +x bin/api

build-playground:
	$(ENV_LOAD); cd playground && \
	  bun run build

build:
	$(MAKE) build-api
	$(MAKE) build-cli

start-api:
	$(ENV_LOAD); ./bin/api

%::
	@true

.PHONY: install install-playground test test-json coverage lint dev run cli dev-cli build-api build-cli build-playground build playground