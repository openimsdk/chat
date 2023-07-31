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

# Use golang as the builder stage
FROM golang:1.20 as builder

WORKDIR /workspace

ENV GOPROXY=https://goproxy.cn

ARG GOARCH
ARG GOOS

# Copy source code files into the image
COPY . .

RUN go mod download

# Compile the source code
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o ./bin/open_im_admin ./cmd/rpc/admin
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o ./bin/open_im_admin_api ./cmd/api/admin_api
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o ./bin/open_im_chat ./cmd/rpc/chat
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o ./bin/open_im_chat_api ./cmd/api/chat_api


# Build the runtime stage
FROM debian

# Set fixed project path
ENV WORKDIR /chat
ENV CMDDIR $WORKDIR/scripts
ENV CONFIG_NAME $WORKDIR/config/config.yaml

# Copy the executable files to the target directory
COPY --from=builder /workspace/bin/open_im_admin $WORKDIR/bin/open_im_admin
COPY --from=builder /workspace/bin/open_im_admin_api $WORKDIR/bin/open_im_admin_api
COPY --from=builder /workspace/bin/open_im_chat $WORKDIR/bin/open_im_chat
COPY --from=builder /workspace/bin/open_im_chat_api $WORKDIR/bin/open_im_chat_api
COPY --from=builder /workspace/scripts $WORKDIR/scripts
COPY --from=builder /workspace/config/config.yaml $WORKDIR/config/config.yaml

# Create several directories for mounting and add executable permissions
RUN mkdir $WORKDIR/logs && \
    chmod +x $WORKDIR/bin/open_im_admin $WORKDIR/bin/open_im_chat $WORKDIR/bin/open_im_admin_api $WORKDIR/bin/open_im_chat_api
RUN apt-get -qq update \
    && apt-get -qq install -y --no-install-recommends ca-certificates curl bash

VOLUME ["/chat/logs","/chat/config","/chat/scripts"]

WORKDIR $CMDDIR
CMD ["./docker_start_all.sh"]