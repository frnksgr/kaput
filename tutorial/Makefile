
image = frnksgr/jupyter-tutorial

.DEFAULT_GOAL := help
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## build image
	docker build -t $(image) .

.PHONY: push
push: ## push image to global registry
	docker push $(image)

.PHONY: run
run: ## run container 
	docker run -it --rm --volume `pwd`/notebooks:/tutorial $(image)

