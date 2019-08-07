TAG := dev-marco-7
.PHONY: examples-docker
examples-docker:
	docker build -t quay.io/solo-io/ext-auth-plugins:$(TAG) .
	docker push quay.io/solo-io/ext-auth-plugins:$(TAG)