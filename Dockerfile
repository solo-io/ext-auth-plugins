# This stage is parametrized to replicate the same environment GlooE was built in.
# All ARGs need to be set via the docker `--build-arg` flags.
ARG GO_BUILD_IMAGE
FROM $GO_BUILD_IMAGE AS build-env

ARG GC_FLAGS
ARG VERIFY_SCRIPT

# Fail if VERIFY_SCRIPT not set
RUN if [[ ! $VERIFY_SCRIPT ]]; then echo "Required VERIFY_SCRIPT build argument not set" && exit 1; fi

RUN apk add --no-cache gcc musl-dev

ADD . /go/src/github.com/solo-io/ext-auth-plugins/
WORKDIR /go/src/github.com/solo-io/ext-auth-plugins

# De-vendor all the dependencies and move them to the GOPATH.
# We need this so that the import paths for any library shared between the plugins and Gloo are the same.
RUN cp -a vendor/. /go/src/ && rm -rf vendor

# Build plugins with CGO enabled
RUN CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -buildmode=plugin -gcflags="$GC_FLAGS" -o examples/RequiredHeader.so examples/required_header/plugin.go

# Verify that plugins can be loaded by GlooE
RUN chmod +x $VERIFY_SCRIPT
RUN $VERIFY_SCRIPT -pluginDir examples -manifest examples/plugin_manifest.yaml

# This stage builds the final image containing just the plugin .so files
FROM alpine:3.10.1
RUN mkdir /compiled-auth-plugins
COPY --from=build-env /go/src/github.com/solo-io/ext-auth-plugins/examples/RequiredHeader.so /compiled-auth-plugins/
CMD cp /compiled-auth-plugins/* /auth-plugins/