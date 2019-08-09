SOURCES=$(shell find . -name "*.go" | grep -v test)

AuthorizeAll.so: $(SOURCES)
	go build -buildmode=plugin -o AuthorizeAll.so plugin.go