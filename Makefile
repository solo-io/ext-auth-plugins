#----------------------------------------------------------------------------------
# Compare dependencies against GlooE
#----------------------------------------------------------------------------------
# Export all variables to sub-makes
export

GLOOE_VERSION=0.18.6
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
# Build and publish example plugin implementations
#----------------------------------------------------------------------------------
publish-example-plugins:
	$(MAKE) -C examples