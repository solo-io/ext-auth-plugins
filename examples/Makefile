#----------------------------------------------------------------------------------
# Build, test and publish example plugins
#----------------------------------------------------------------------------------
REPO_ROOT=./..
SOURCES=$(shell find . -name "*.go" | grep -v test)
TAG := dev

.PHONY: publish-examples
publish-examples: test-examples
	cd $(REPO_ROOT) && docker build -t quay.io/solo-io/ext-auth-plugins:$(TAG) .
	cd $(REPO_ROOT) && docker push quay.io/solo-io/ext-auth-plugins:$(TAG)

.PHONY: test-examples
test-examples: build-examples
	ginkgo ./...

.PHONY: build-examples
build-examples: $(SOURCES)
	go build -buildmode=plugin -o authorize_all/AuthorizeAll.so authorize_all/plugin.go
	go build -buildmode=plugin -o header/RequiredHeader.so header/plugin.go