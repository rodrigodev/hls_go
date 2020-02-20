FROM golang:1.13 as builder

RUN mkdir -p /hls_go/

WORKDIR /hls_go

COPY . .

RUN go mod download

RUN go test -v -race ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -o bin/hls_go pkg/*

FROM alpine:3.10

RUN addgroup -S app \
    && adduser -S -g app app \
    && apk --no-cache add \
    curl openssl netcat-openbsd

WORKDIR /home/app

COPY --from=builder /hls_go/bin/hls_go .
COPY --from=builder /hls_go/bin/hls_go /usr/local/bin/hls_go
RUN chown -R app:app ./

USER app

CMD ["./hls_go"]
