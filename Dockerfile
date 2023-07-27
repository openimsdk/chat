# Use the official Golang image as a build environment
# docker buildx build --platform linux/amd64,linux/arm64 -t  ghcr.io/openimsdk/openim-chat:v1.0.1 . --push

FROM golang:1.20 AS builder

WORKDIR /app

ARG GOARCH
ARG GOOS

ENV GOPROXY=https://goproxy.cn

# Copy go mod and go sum files then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code files into the image
COPY . .

# Compile the source code
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o open_im_admin ./cmd/rpc/admin
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o open_im_admin_api ./cmd/api/admin_api
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o open_im_chat ./cmd/rpc/chat
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o open_im_chat_api ./cmd/api/chat_api

# Create a new image layer using scratch, copy the built binary files into it
FROM alpine:latest

# Create some directories for mounting, add execution permissions
RUN mkdir -p $WORKDIR/logs \
    && chmod +x $WORKDIR/bin/open_im_admin $WORKDIR/bin/open_im_chat $WORKDIR/bin/open_im_admin_api $WORKDIR/bin/open_im_chat_api \
    && apk --no-cache add ca-certificates curl 

ENV WORKDIR /chat
ENV CMDDIR $WORKDIR/scripts
ENV CONFIG_NAME $WORKDIR/config/config.yaml

COPY ./scripts $WORKDIR/scripts
COPY ./config/config.yaml $WORKDIR/config/config.yaml
COPY --from=builder /app/open_im_admin $WORKDIR/bin/open_im_admin
COPY --from=builder /app/open_im_admin_api $WORKDIR/bin/open_im_admin_api
COPY --from=builder /app/open_im_chat $WORKDIR/bin/open_im_chat
COPY --from=builder /app/open_im_chat_api $WORKDIR/bin/open_im_chat_api

VOLUME ["/chat/logs","/chat/config","/chat/scripts"]

WORKDIR $CMDDIR
CMD ["./docker_start_all.sh"]