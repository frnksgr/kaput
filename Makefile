
IMAGE 		= gcr.io/sap-cp-gke/kaput
CF_IMAGE 	= frnksgr/kaput

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

.PHONY: docker-cf-build
docker-cf-build: ## build docker image for CF
	docker build -t $(CF_IMAGE) --build-arg BASEIMAGE=alpine:3.9 .

.PHONY: docker-cf-push
docker-cf-push: docker-cf-build ## push CF docker-image
	docker push $(CF_IMAGE)

