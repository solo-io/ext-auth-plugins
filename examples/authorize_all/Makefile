SOURCES=$(shell find . -name "*.go" | grep -v test)

AuthorizeAll.so: $(SOURCES)
	CGO_ENABLED=1 GOARCH=amd64 GOOS=darwin go build -buildmode=plugin -gcflags='all=-N -l' -o AuthorizeAll.so plugin.go