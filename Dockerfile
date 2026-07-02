FROM golang:alpine AS builder

ARG RELEASE=false
ARG COMPRESS=false
WORKDIR /openim-chat

RUN apk add --no-cache upx
RUN go install github.com/magefile/mage@latest

COPY . .
RUN go mod download
RUN RELEASE=${RELEASE} COMPRESS=${COMPRESS} mage build
RUN mage -compile ./mage -ldflags "-s -w"

FROM alpine:latest

WORKDIR /openim-chat

RUN apk add --no-cache bash

COPY --from=builder /openim-chat/_output ./_output
COPY --from=builder /openim-chat/config ./config
COPY --from=builder /openim-chat/start-config.yml ./start-config.yml
COPY --from=builder /openim-chat/mage ./mage

ENTRYPOINT ["sh", "-c", "./mage start && sleep infinity"]
