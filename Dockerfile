FROM golang:1.16-alpine as build


LABEL Maintainer="Rodrigo Carneiro <rodrigo.carneiro.dev@gmail.com>"

ARG COMMIT='local'
ARG TIMESTAMP='local'

WORKDIR /go/src/github.com/rodrigodev/hls_go

COPY . .

RUN go mod vendor

RUN apk update add ca-certificates && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod vendor -ldflags "-s -w -X 'main.githash=${COMMIT}' -X 'main.buildstamp=${TIMESTAMP}'" -o /app ./pkg

FROM scratch
#FROM alpine:3.4

#/go/src/github.com/rodrigodev/hls_go
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go/src/github.com/rodrigodev/hls_go/static static
COPY --from=build /go/src/github.com/rodrigodev/hls_go/media media
COPY --from=build /app app

EXPOSE 8080

ENTRYPOINT ["./app"]
