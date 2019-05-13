
IMAGE = gcr.io/sap-cp-gke/kaput

.DEFAULT_GOAL := help
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: docker-build
docker-build: ## build docker image
	docker build -t $(IMAGE) .

.PHONY: docker-push
docker-push: docker-build ## push docker-image
	docker push $(IMAGE)

