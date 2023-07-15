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

# Use an existing docker image as base
FROM ubuntu

# Set work directory
WORKDIR /chat

# Copy files from project to the work directory
COPY ./config /chat/config
COPY ./scripts /chat/scripts
COPY ./logs /chat/logs
COPY ./bin /chat/bin

# Make the script executable
RUN chmod +x ./scripts/docker_start_all.sh

# Create volumes for these directories
VOLUME ["/chat/logs"]

WORKDIR /chat/scripts


# Run the script when the container starts
CMD ["./docker_start_all.sh"]
