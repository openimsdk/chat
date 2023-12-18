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
FROM golang:1.21 AS builder

ARG GO111MODULE=on
ARG GOPROXY=https://goproxy.cn,direct

WORKDIR /openim/openim-chat

ENV GO111MODULE=$GO111MODULE
ENV GOPROXY=$GOPROXY

COPY go.mod go.sum ./
RUN go mod download

# Copy all files to the container
ADD . .

RUN make clean
RUN make build

# Build the runtime stage
FROM ghcr.io/openim-sigs/openim-ubuntu-image:latest

WORKDIR ${CHAT_WORKDIR}

COPY --from=builder ${OPENIM_CHAT_BINDIR} /openim/openim-chat/_output/bin
COPY --from=builder ${CHAT_WORKDIR}/config /openim/openim-chat/config
COPY --from=builder ${CHAT_WORKDIR}/scripts /openim/openim-chat/scripts
COPY --from=builder ${CHAT_WORKDIR}/deployments /openim/openim-chat/deployments

CMD ["/openim/openim-chat/scripts/docker_start_all.sh"]
