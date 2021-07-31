FROM golang:1.16.5-alpine3.13 AS builder

RUN set -eux; \
        apk update && \
        apk add --no-cache --virtual .build-deps gcc libc-dev git

ENV GOPATH /go/
ENV GO_WORKDIR $GOPATH/src/raft-app
ENV GO111MODULE=on

WORKDIR $GO_WORKDIR
ADD . $GO_WORKDIR
# add commit num into binary for checking
RUN go build -ldflags "-X main.gitcommitnum=`git log|head -1|awk '{print $2}'` -s -w" -o raft-app


FROM alpine:3.13.5
WORKDIR /app
COPY --from=builder /go/src/raft-app/raft-app .
# add app.dev.ini into Dockerfile
RUN mkdir conf
COPY ./conf/app.dev.ini ./conf/

CMD ./raft-app
