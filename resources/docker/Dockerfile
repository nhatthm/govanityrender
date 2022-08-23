FROM golang:1.19 AS builder

ARG VERSION=dev
ARG REVISION=""
ARG USER=""

ENV VERSION=$VERSION
ENV REVISION=$REVISION
ENV USER=$USER

WORKDIR /app/

COPY ../.. .

RUN GOARCH=amd64 GOOS=linux CGO_ENABLED=0 make build

FROM alpine:latest AS final

WORKDIR /tmp

RUN apk add --no-cache git

COPY --from=builder /app/out/vanityrender /bin/vanityrender

CMD ["/bin/vanityrender"]
