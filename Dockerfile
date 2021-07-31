FROM golang:1.16.5-alpine3.13 AS builder

RUN set -eux; \
        apk update && \
        apk add --no-cache --virtual .build-deps gcc libc-dev git

ENV GOPATH /go/
ENV GO_WORKDIR $GOPATH/src/go-third-party/raft-apps/app
ENV GO111MODULE=on

WORKDIR $GO_WORKDIR
ADD . $GO_WORKDIR
RUN go build 

FROM alpine:3.13.5
WORKDIR /app
COPY --from=builder /go/src/go-third-party/raft-apps/app/app .
# TODO : add config.ini into Dockerfile
RUN mkdir conf
COPY ./conf/app.dev.ini ./conf/

CMD ./app
