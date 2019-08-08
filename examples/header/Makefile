SOURCES=$(shell find . -name "*.go" | grep -v test)

RequiredHeader.so: $(SOURCES)
	go build -buildmode=plugin -o RequiredHeader.so plugin.go