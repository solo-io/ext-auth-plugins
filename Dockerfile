FROM golang:alpine AS build-env
RUN apk add --no-cache gcc musl-dev git
ADD . /go/src/github.com/solo-io/ext-auth-plugins/
WORKDIR /go/src/github.com/solo-io/ext-auth-plugins

RUN go get ./example/...
RUN go get golang.org/x/net/trace

RUN rm -rf /go/src/github.com/solo-io/ext-auth-plugins/vendor/golang.org/x/net/trace
RUN ls -lah /go/src/github.com/solo-io/ext-auth-plugins/vendor/golang.org/x/net

RUN CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -buildmode=plugin -gcflags='all=-N -l' -o AuthorizeAll.so example/authorize-all/plugin.go
RUN CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -buildmode=plugin -gcflags='all=-N -l' -o RequiredHeader.so example/header/plugin.go

FROM alpine:3.10.1
RUN mkdir /compiled-auth-plugins
COPY --from=build-env /go/src/github.com/solo-io/ext-auth-plugins/AuthorizeAll.so /compiled-auth-plugins/
COPY --from=build-env /go/src/github.com/solo-io/ext-auth-plugins/RequiredHeader.so /compiled-auth-plugins/
CMD cp /compiled-auth-plugins/* /auth-plugins/