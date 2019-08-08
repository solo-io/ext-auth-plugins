#----------------------------------------------------------------------------------
# Compare dependencies against GlooE
#----------------------------------------------------------------------------------
GLOOE_VERSION=0.18.5

.PHONY: compare-deps
compare-deps: GlooE-Gopkg.lock
	go run scripts/compare_dependencies.go

GlooE-Gopkg.lock:
	curl -o GlooE-Gopkg.lock http://storage.googleapis.com/gloo-ee-dependencies/$(GLOOE_VERSION)/Gopkg.lock


#----------------------------------------------------------------------------------
# Build, test and publish example plugins
#----------------------------------------------------------------------------------
TAG := dev

.PHONY: publish-examples
publish-examples:
	docker build -t quay.io/solo-io/ext-auth-plugins:$(TAG) .
	docker push quay.io/solo-io/ext-auth-plugins:$(TAG)

.PHONY: test-examples
test-examples: subsystem
	cd examples && ginkgo ./...

.PHONY: subsystem
subsystem:
	$(MAKE) -C examples/authorize_all
	$(MAKE) -C examples/header