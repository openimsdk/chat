# Copyright © 2023 OpenIM open source community. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

ARG GOARCH
ARG GOOS

# Use golang as the builder stage
FROM golang:1.20 AS builder

ARG GO111MODULE=on
ARG GOPROXY=https://goproxy.cn,direct

WORKDIR /openim/openim-chat

ENV GO111MODULE=$GO111MODULE
ENV GOPROXY=$GOPROXY

COPY . .

RUN go mod download

# Compile the source code
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o /openim/openim-chat/bin/open_im_admin ./cmd/rpc/admin
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o /openim/openim-chat/bin/open_im_admin_api ./cmd/api/admin_api
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o /openim/openim-chat/bin/open_im_chat ./cmd/rpc/chat
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o /openim/openim-chat/bin/open_im_chat_api ./cmd/api/chat_api

# Build the runtime stage
FROM ghcr.io/openim-sigs/openim-bash-image:v1.3.0

# Set fixed project path
WORKDIR /openim/openim-chat

# Copy the executable files to the target directory
COPY --from=builder ${OPENIM_CHAT_BINDIR} /openim/openim-chat/bin

COPY --from=builder ${OPENIM_CHAT_CMDDIR} /openim/openim-chat/scripts
COPY --from=builder ${OPENIM_CHAT_CONFIG_NAME} /openim/openim-chat/config/config.yaml

WORKDIR $OPENIM_CHAT_CMDDIR

CMD ["./docker_start_all.sh"]