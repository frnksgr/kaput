

.DEFAULT_GOAL := help

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## build local to create bin/... executable(s)
	go build -v ./pkg/...
	go build -o ./bin/kaput-server ./cmd/kaput-server/...


.PHONY: clean
clean: ## clean up
	go clean -i ./...
	rm -f bin/kaput-server
