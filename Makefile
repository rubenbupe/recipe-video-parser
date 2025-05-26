# Cargar variables de entorno desde .env si existe
ENV_LOAD = [ -f .env ] && set -a && . ./.env && set +a || true


ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

# Dependencies
install:
	go get ./...


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

build-api:
	CGO_ENABLED=0 go build -o bin/api cmd/api/main.go

build-playground:
	$(ENV_LOAD); cd playground && \
	  bun run build

build:
	$(MAKE) build-api
	$(MAKE) build-cli
	$(MAKE) build-playground

%::
	@true

.PHONY: install test test-json coverage lint dev run cli dev-cli build-api build-cli build-playground build playground