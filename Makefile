GLOOE_VERSION := 0.18.10
BUILD_ID := $(BUILD_ID)
RELEASE := "true"
ifeq ($(TAGGED_VERSION),)
	TAGGED_VERSION := v$(BUILD_ID)
	RELEASE := "false"
endif
VERSION ?= $(shell echo $(TAGGED_VERSION) | cut -c 2-)

#----------------------------------------------------------------------------------
# Retrieve GlooE build information
#----------------------------------------------------------------------------------
GLOOE_DIR := _glooe
z := $(shell mkdir -p $(GLOOE_DIR))

.PHONY: get-glooe-info
get-glooe-info: $(GLOOE_DIR)/Gopkg.lock $(GLOOE_DIR)/verify-plugins-linux-amd64 $(GLOOE_DIR)/build_env

$(GLOOE_DIR)/Gopkg.lock:
	curl -o $@ http://storage.googleapis.com/gloo-ee-dependencies/$(GLOOE_VERSION)/Gopkg.lock

$(GLOOE_DIR)/verify-plugins-linux-amd64:
	curl -o $@ http://storage.googleapis.com/gloo-ee-dependencies/$(GLOOE_VERSION)/verify-plugins-linux-amd64

$(GLOOE_DIR)/build_env:
	curl -o $@ http://storage.googleapis.com/gloo-ee-dependencies/$(GLOOE_VERSION)/build_env


#----------------------------------------------------------------------------------
# Compare dependencies against GlooE
#----------------------------------------------------------------------------------

.PHONY: compare-deps
compare-deps: Gopkg.lock $(GLOOE_DIR)/Gopkg.lock
	go run scripts/compare_dependencies.go Gopkg.lock $(GLOOE_DIR)/Gopkg.lock


#----------------------------------------------------------------------------------
# Build plugins
#----------------------------------------------------------------------------------
EXAMPLES_DIR := examples
SOURCES := $(shell find . -name "*.go" | grep -v test)

define get_glooe_var
$(shell grep $(1) $(GLOOE_DIR)/build_env | cut -d '=' -f2)
endef

.PHONY: build-plugins
build-plugins: $(GLOOE_DIR)/build_env $(GLOOE_DIR)/verify-plugins-linux-amd64
	docker build -t quay.io/solo-io/ext-auth-plugins:$(VERSION) \
		--build-arg GO_BUILD_IMAGE=$(call get_glooe_var,GO_BUILD_IMAGE) \
		--build-arg GC_FLAGS=$(call get_glooe_var,GC_FLAGS) \
		--build-arg VERIFY_SCRIPT=$(GLOOE_DIR)/verify-plugins-linux-amd64 \
		.

.PHONY: build-plugins-for-tests
build-plugins-for-tests: $(EXAMPLES_DIR)/header/RequiredHeader.so

$(EXAMPLES_DIR)/header/RequiredHeader.so: $(SOURCES)
	go build -buildmode=plugin -o $(EXAMPLES_DIR)/header/RequiredHeader.so $(EXAMPLES_DIR)/header/plugin.go


#----------------------------------------------------------------------------------
# Release plugins
#----------------------------------------------------------------------------------

.PHONY: release-plugins
release-plugins: build-plugins
ifeq ($(RELEASE),"true")
	docker push quay.io/solo-io/ext-auth-plugins:$(VERSION)
else
	@echo This is not a release build. Example plugins will not be published.
endif