FROM golang:alpine AS build-env
RUN apk add --no-cache gcc musl-dev
ADD . /go/src/github.com/solo-io/ext-auth-plugins/
WORKDIR /go/src/github.com/solo-io/ext-auth-plugins
RUN go get ./example/...
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -buildmode=plugin -o AuthorizeAll.so example/authorize-all/plugin.go
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -buildmode=plugin -o RequiredHeader.so example/header/plugin.go

FROM alpine
RUN mkdir /compiled-auth-plugins
COPY --from=build-env /go/src/github.com/solo-io/ext-auth-plugins/AuthorizeAll.so /compiled-auth-plugins/
COPY --from=build-env /go/src/github.com/solo-io/ext-auth-plugins/RequiredHeader.so /compiled-auth-plugins/
CMD cp /compiled-auth-plugins/* /auth-plugins/