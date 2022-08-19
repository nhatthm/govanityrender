FROM golang:1.19 AS builder

WORKDIR /app/

COPY ../.. .

RUN GOARCH=amd64 GOOS=linux CGO_ENABLED=0 make build

FROM alpine:latest AS final

WORKDIR /tmp

RUN apk add --no-cache git

COPY --from=builder /app/out/vanityrender /bin/vanityrender

CMD ["/bin/vanityrender"]
