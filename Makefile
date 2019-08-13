#----------------------------------------------------------------------------------
# Compare dependencies against GlooE
#----------------------------------------------------------------------------------
# Export all variables to sub-makes
export

GLOOE_VERSION := 0.18.6
BUILD_ID := $(BUILD_ID)
RELEASE := "true"
ifeq ($(TAGGED_VERSION),)
	TAGGED_VERSION := v$(BUILD_ID)
	RELEASE := "false"
endif
VERSION ?= $(shell echo $(TAGGED_VERSION) | cut -c 2-)

.PHONY: compare-deps
compare-deps: Gopkg.lock GlooE-Gopkg.lock print-info
	go run scripts/compare_dependencies.go Gopkg.lock GlooE-Gopkg.lock

GlooE-Gopkg.lock:
	curl -o GlooE-Gopkg.lock http://storage.googleapis.com/gloo-ee-dependencies/$(GLOOE_VERSION)/Gopkg.lock

# TODO: remove
.PHONY: print-info
print-info:
	@echo BUILD_ID: $(BUILD_ID)
	@echo TAGGED_VERSION: $(TAGGED_VERSION)
	@echo VERSION: $(VERSION)
	@echo RELEASE: $(RELEASE)

#----------------------------------------------------------------------------------
# Build, test and publish example plugins
#----------------------------------------------------------------------------------
EXAMPLES_DIR := examples
SOURCES := $(shell find . -name "*.go" | grep -v test)

.PHONY: publish-examples
publish-examples:
ifeq ($(RELEASE),"true")
	docker build -t quay.io/solo-io/ext-auth-plugins:$(VERSION) .
	docker push quay.io/solo-io/ext-auth-plugins:$(VERSION)
else
	@echo This is not a release build. Example plugins will not be published.
endif

.PHONY: build-examples-for-tests
build-examples-for-tests: $(EXAMPLES_DIR)/authorize_all/AuthorizeAll.so $(EXAMPLES_DIR)/header/RequiredHeader.so

$(EXAMPLES_DIR)/authorize_all/AuthorizeAll.so: $(SOURCES)
	go build -buildmode=plugin -o $(EXAMPLES_DIR)/authorize_all/AuthorizeAll.so $(EXAMPLES_DIR)/authorize_all/plugin.go

$(EXAMPLES_DIR)/header/RequiredHeader.so: $(SOURCES)
	go build -buildmode=plugin -o $(EXAMPLES_DIR)/header/RequiredHeader.so $(EXAMPLES_DIR)/header/plugin.go