FROM golang:1.12.7-alpine AS build-env
RUN apk add --no-cache gcc musl-dev git

ADD . /go/src/github.com/solo-io/ext-auth-plugins/
WORKDIR /go/src/github.com/solo-io/ext-auth-plugins

# TODO(marco): check if we really need this
RUN go get ./examples/...

# Go get `golang.org/x/net/trace` and remove it from vendor.
# If this package is initialized more than once, it causes a panic (see issue: https://github.com/golang/go/issues/24137).
# We want this to reside ONLY in the GOPATH, both for the ext-auth server and any go plugins it has to load. This way
# the import paths for the package will be the same and the `init()` function will be called only once.
RUN go get golang.org/x/net/trace
RUN rm -rf /go/src/github.com/solo-io/ext-auth-plugins/vendor/golang.org/x/net/trace

# Build plugins with CGO enabled
RUN CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -buildmode=plugin -gcflags='all=-N -l' -o AuthorizeAll.so examples/authorize_all/plugin.go
RUN CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -buildmode=plugin -gcflags='all=-N -l' -o RequiredHeader.so examples/header/plugin.go

FROM alpine:3.10.1
RUN mkdir /compiled-auth-plugins
COPY --from=build-env /go/src/github.com/solo-io/ext-auth-plugins/AuthorizeAll.so /compiled-auth-plugins/
COPY --from=build-env /go/src/github.com/solo-io/ext-auth-plugins/RequiredHeader.so /compiled-auth-plugins/
CMD cp /compiled-auth-plugins/* /auth-plugins/