PROJECT_BINARY=sidepeer
PROJECT_BINARY_OUTPUT=out
PROJECT_RELEASER_OUTPUT=dist

.PHONY: all

all: help

## Build:
tidy: ## Tidy project
	go mod tidy

clean: ## Cleans temporary folder
	rm -rf /tmp/peer*
	rm -rf ${PROJECT_BINARY_OUTPUT}
	rm -rf ${PROJECT_RELEASER_OUTPUT}

build: clean tidy ## Builds project
	GO111MODULE=on CGO_ENABLED=0 go build -ldflags="-w -s" -o ${PROJECT_BINARY_OUTPUT}/bin/${PROJECT_BINARY} cmd/${PROJECT_BINARY}/main.go

test: clean tidy ## Run unit tests
	go test -v

coverage: clean tidy ## Run code coverage
	go test -cover	

## GoReleaser:
release-local: clean tidy ## Creates snapshot releases
	goreleaser release --snapshot --rm-dist

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    %-20s%s\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  %s\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
