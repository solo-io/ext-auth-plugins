FROM golang:alpine AS build-env
RUN apk add --no-cache gcc musl-dev
ADD . /go/src/github.com/solo-io/ext-auth-plugins/
WORKDIR /go/src/github.com/solo-io/ext-auth-plugins
RUN go get ./example/...
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -buildmode=plugin -o AuthorizeAll.so example/authorize-all/plugin.go

FROM alpine
RUN mkdir /auth-plugins
COPY --from=build-env /go/src/github.com/solo-io/ext-auth-plugins/AuthorizeAll.so /auth-plugins/