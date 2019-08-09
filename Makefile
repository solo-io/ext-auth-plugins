#----------------------------------------------------------------------------------
# Compare dependencies against GlooE
#----------------------------------------------------------------------------------
GLOOE_VERSION=dev

.PHONY: compare-deps
compare-deps: Gopkg.lock GlooE-Gopkg.lock
	go run scripts/compare_dependencies.go Gopkg.lock GlooE-Gopkg.lock

GlooE-Gopkg.lock:
	curl -o GlooE-Gopkg.lock http://storage.googleapis.com/gloo-ee-dependencies/$(GLOOE_VERSION)/Gopkg.lock