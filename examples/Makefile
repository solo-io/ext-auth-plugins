#----------------------------------------------------------------------------------
# Build, test and publish example plugins
#----------------------------------------------------------------------------------
REPO_ROOT=./..
SOURCES=$(shell find . -name "*.go" | grep -v test)

.PHONY: publish-examples
publish-examples:
	cd $(REPO_ROOT) && docker build -t quay.io/solo-io/ext-auth-plugins:$(VERSION) .
	docker push quay.io/solo-io/ext-auth-plugins:$(VERSION)

.PHONY: test-examples
test-examples: build-examples-for-tests
	ginkgo ./...

.PHONY: build-examples-for-tests
build-examples-for-tests: $(SOURCES)
	go build -buildmode=plugin -o authorize_all/AuthorizeAll.so authorize_all/plugin.go
	go build -buildmode=plugin -o header/RequiredHeader.so header/plugin.go