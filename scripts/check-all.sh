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





DATA="$(date +%H:%M:%S)"
echo "# Start Chat check-all.sh ${DATA}, For local deployments, use ./check-all.sh --print-screen"

# 检查第一个参数是否为 --print-screen
if [ "$1" == "--print-screen" ]; then
    PRINT_SCREEN=1
fi

mkdir -p ${SCRIPTS_ROOT}/../logs

if [ -z "$PRINT_SCREEN" ]; then
    exec >> ${SCRIPTS_ROOT}/../logs/chat_$(date '+%Y%m%d').log 2>&1
fi

#Include shell font styles and some basic information
source $SCRIPTS_ROOT/style-info.sh
source $SCRIPTS_ROOT/path-info.sh
source $SCRIPTS_ROOT/function.sh
source $SCRIPTS_ROOT/util.sh


all_services_running=true

for binary_path in "${binary_full_paths[@]}"; do
    check_services_with_name "$binary_path"
    if [ $? -ne 0 ]; then
        all_services_running=false
        # Print the binary path in red for not running services
        echo -e "\033[0;31mService not running: $binary_path\033[0m"
    fi
done

if $all_services_running; then
    # Print "Startup successful" in green
    echo -e "\033[0;32mStartup successful\033[0m"
else
    echo "One or more services are not running."
fi


