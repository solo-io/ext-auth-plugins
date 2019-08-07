
.PHONY: examples-docker
examples-docker:
	docker build -t quay.io/solo-io/ext-auth-plugins:dev-marco-2 .
	docker push quay.io/solo-io/ext-auth-plugins:dev-marco-2