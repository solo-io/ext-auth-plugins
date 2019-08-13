FROM golang:1.12.7-alpine AS build-env

RUN apk add --no-cache gcc musl-dev git

ADD . /go/src/github.com/solo-io/ext-auth-plugins/
WORKDIR /go/src/github.com/solo-io/ext-auth-plugins

# De-vendor all the dependencies and move them to the GOPATH.
# We need this so that the import paths for any library shared between the plugins and Gloo are the same.
RUN cp -a vendor/. /go/src/ && rm -rf vendor

# Build plugins with CGO enabled
RUN CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -buildmode=plugin -gcflags='all=-N -l' -o AuthorizeAll.so examples/authorize_all/plugin.go
RUN CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -buildmode=plugin -gcflags='all=-N -l' -o RequiredHeader.so examples/header/plugin.go

FROM alpine:3.10.1
RUN mkdir /compiled-auth-plugins
COPY --from=build-env /go/src/github.com/solo-io/ext-auth-plugins/AuthorizeAll.so /compiled-auth-plugins/
COPY --from=build-env /go/src/github.com/solo-io/ext-auth-plugins/RequiredHeader.so /compiled-auth-plugins/
CMD cp /compiled-auth-plugins/* /auth-plugins/