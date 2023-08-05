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

COPY go.mod go.sum ./
RUN go mod download

# Copy all files to the container
ADD . .

RUN /bin/sh -c "make clean"
RUN /bin/sh -c "make build"

# Build the runtime stage
FROM ghcr.io/openim-sigs/openim-bash-image:latest

WORKDIR ${CHAT_WORKDIR}

COPY --from=builder /openim/openim-chat/_output/bin/platforms /openim/openim-chat/_output/bin/platforms
COPY --from=builder ${OPENIM_CHAT_CMDDIR} /openim/openim-chat/scripts
COPY --from=builder ${OPENIM_CHAT_CONFIG_NAME} /openim/openim-chat/config/config.yaml

VOLUME ["/openim/openim-chat/_output","/openim/openim-chat/logs","/openim/openim-chat/config","/openim/openim-chat/scripts"]

CMD ${OPENIM_CHAT_CMDDIR}/docker_start_all.sh
