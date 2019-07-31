FROM golang:1.12-alpine3.10

RUN apk add --no-cache \
    ca-certificates \
    gcc \
    libc-dev \
    git

ENV GO111MODULE=on
EXPOSE 8002
WORKDIR /work/

ENTRYPOINT ["./docker-entrypoint.sh"]
