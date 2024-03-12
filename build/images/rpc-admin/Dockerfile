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

# OpenIM base image: https://github.com/openim-sigs/openim-base-image

# Set go mod installation source and proxy

ARG GOARCH
ARG GOOS

FROM golang:1.21 AS builder

ARG GO111MODULE=on
ARG GOPROXY=https://goproxy.io,direct

WORKDIR /openim/openim-chat

ENV GO111MODULE=$GO111MODULE
ENV GOPROXY=$GOPROXY

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build BINS=admin-rpc

RUN cp /openim/openim-chat/_output/bin/platforms/$(go env GOOS)/$(go env GOARCH)/admin-rpc /usr/bin/admin-rpc


FROM ghcr.io/openim-sigs/openim-bash-image:latest

WORKDIR ${CHAT_WORKDIR}

COPY --from=builder /usr/bin/admin-rpc ./bin/admin-rpc

ENTRYPOINT ["./bin/admin-rpc"]