
IMAGE 		= docker.io/frnksgr/kaput

.DEFAULT_GOAL := help
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: docker-build
docker-build: ## build docker image
	docker build -t $(IMAGE) .
	docker tag $(IMAGE) $(IMAGE):scratch
	docker build -t $(IMAGE):alpine3.9 --build-arg BASEIMAGE=alpine:3.9 .

.PHONY: docker-push
docker-push: docker-build ## push docker-image
	docker push $(IMAGE)
	docker push $(IMAGE):scratch
	docker push $(IMAGE):alpine3.9
