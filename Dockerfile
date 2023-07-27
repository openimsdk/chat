# 使用官方 Golang 镜像作为构建环境
FROM golang:1.20 AS builder

WORKDIR /app

ARG GOARCH
ARG GOOS

ENV GOPROXY=https://goproxy.cn

# 复制 go mod 和 go sum 文件然后下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源码文件到镜像中
COPY . .

# 编译源码
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o open_im_admin ./cmd/rpc/admin
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o open_im_admin_api ./cmd/api/admin_api
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o open_im_chat ./cmd/rpc/chat
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o open_im_chat_api ./cmd/api/chat_api

# 使用 scratch 创建一个新的镜像层，复制构建好的二进制文件进去
FROM ubuntu

ENV WORKDIR /chat
ENV CMDDIR $WORKDIR/scripts
ENV CONFIG_NAME $WORKDIR/config/config.yaml

COPY --from=builder /app/open_im_admin $WORKDIR/bin/open_im_admin
COPY --from=builder /app/open_im_admin_api $WORKDIR/bin/open_im_admin_api
COPY --from=builder /app/open_im_chat $WORKDIR/bin/open_im_chat
COPY --from=builder /app/open_im_chat_api $WORKDIR/bin/open_im_chat_api
COPY ./scripts $WORKDIR/scripts
COPY ./config/config.yaml $WORKDIR/config/config.yaml

# 创建用于挂载的几个目录，添加可执行权限
RUN mkdir $WORKDIR/logs && \
    chmod +x $WORKDIR/bin/open_im_admin $WORKDIR/bin/open_im_chat $WORKDIR/bin/open_im_admin_api $WORKDIR/bin/open_im_chat_api
RUN apt-get -qq update \
    && apt-get -qq install -y --no-install-recommends ca-certificates curl

VOLUME ["/chat/logs","/chat/config","/chat/scripts"]

WORKDIR $CMDDIR
CMD ["./docker_start_all.sh"]
