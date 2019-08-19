ARG GO_BUILD_IMAGE
FROM $GO_BUILD_IMAGE AS build-env

# Get VERIFY_SCRIPT args and fail if not set
ARG VERIFY_SCRIPT
RUN if [[ ! $VERIFY_SCRIPT ]]; then echo "Required VERIFY_SCRIPT build argument not set" && exit 1; fi

RUN apk add --no-cache gcc musl-dev

ADD . /go/src/github.com/solo-io/ext-auth-plugins/
WORKDIR /go/src/github.com/solo-io/ext-auth-plugins

# De-vendor all the dependencies and move them to the GOPATH.
# We need this so that the import paths for any library shared between the plugins and Gloo are the same.
RUN cp -a vendor/. /go/src/ && rm -rf vendor

# Build plugins with CGO enabled
RUN CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -buildmode=plugin -gcflags='all=-N -l' \
    -o examples/compiled/AuthorizeAll.so examples/authorize_all/plugin.go
RUN CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -buildmode=plugin -gcflags='all=-N -l' \
    -o examples/compiled/RequiredHeader.so examples/header/plugin.go

# Verify that plugins can be loaded by GlooE
RUN chmod +x $VERIFY_SCRIPT
RUN $VERIFY_SCRIPT -pluginDir examples/compiled -f examples/plugin_manifest.yaml -debug

FROM alpine:3.10.1
RUN mkdir /compiled-auth-plugins
COPY --from=build-env /go/src/github.com/solo-io/ext-auth-plugins/examples/compiled/AuthorizeAll.so /compiled-auth-plugins/
COPY --from=build-env /go/src/github.com/solo-io/ext-auth-plugins/examples/compiled/RequiredHeader.so /compiled-auth-plugins/
CMD cp /compiled-auth-plugins/* /auth-plugins/