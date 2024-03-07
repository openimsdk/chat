#!/usr/bin/env bash


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

SCRIPTS_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
OPENIM_ROOT=$(dirname "${SCRIPTS_ROOT}")/..
source $SCRIPTS_ROOT/util.sh


logs_dir="$SCRIPTS_ROOT/../_output/logs"
DOCKER_LOG_FILE="$logs_dir/chat-docker.log"
echo 111111111111111111
if is_running_in_container; then
  exec >> ${DOCKER_LOG_FILE} 2>&1
fi
echo 111111111111111111


DATA="$(date +%H:%M:%S)"
echo "# Start Chat check-all.sh ${DATA}, For local deployments, use ./check-all.sh --print-screen"

# --print-screen
if [ "$1" == "--print-screen" ]; then
    PRINT_SCREEN=1
fi

#mkdir -p ${SCRIPTS_ROOT}/../logs

#if [ -z "$PRINT_SCREEN" ]; then
#    exec >> ${SCRIPTS_ROOT}/../logs/chat_$(date '+%Y%m%d').log 2>&1
#fi

#Include shell font styles and some basic information
source $SCRIPTS_ROOT/style-info.sh
source $SCRIPTS_ROOT/path-info.sh
source $SCRIPTS_ROOT/function.sh
source $SCRIPTS_ROOT/util.sh


all_services_running=true
not_running_count=0 # Initialize a counter for not running services

for binary_path in "${binary_full_paths[@]}"; do
    result=$(check_services_with_name "$binary_path")
    if [ $? -ne 0 ]; then
        all_services_running=false
        not_running_count=$((not_running_count + 1)) # Increment the counter
        # Print the binary path in red for not running services
        echo -e "\033[0;31mService not running: $binary_path\033[0m"
    fi
    exit 1
done

if $all_services_running; then
    # Print "Startup successful" in green
    echo -e "\033[0;32mAll chat services startup successful\033[0m"
else
    # Print the number of services that are not running
    echo -e "\033[0;31m$not_running_count chat service(s) are not running.\033[0m"
    exit 1
fi
