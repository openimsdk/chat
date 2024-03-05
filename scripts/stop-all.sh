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

#fixme This scripts is to stop the service
SCRIPTS_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
OPENIM_ROOT=$(cd $(dirname "${BASH_SOURCE[0]}")/.. &&pwd)

source $OPENIM_ROOT/scripts/style-info.sh
source $OPENIM_ROOT/scripts/path-info.sh
source $SCRIPTS_ROOT/function.sh
source $SCRIPTS_ROOT/util.sh



# Loop through each binary full path and attempt to stop the service
for binary_path in "${binary_full_paths[@]}"; do
  result=$(stop_services_with_name "$binary_path")
    ret_val=$?
    if [ $ret_val -ne 0 ]; then
        # Print detailed error log if stop_services_with_name function returns a non-zero value
        echo "Error stopping service at path $binary_path"
    fi
done


all_services_stopped=true

for binary_path in "${binary_full_paths[@]}"; do
  result=$(check_services_with_name "$binary_path")
    if [ $? -eq 0 ]; then
        all_services_stopped=false
        # Print the binary path in red to indicate the service is still running
        echo -e "\033[0;31mService still running: $binary_path\033[0m"
    fi
done

if $all_services_stopped; then
    # Print "All services stopped" in green to indicate success
    echo -e "\033[0;32mAll chat services stopped\033[0m"
else
    # Print error message indicating not all services are stopped
    echo -e "\033[0;31mError: Not all chat services have been stopped.\033[0m"
fi


