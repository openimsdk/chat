#!/usr/bin/env bash
# Copyright © 2023 OpenIM. All rights reserved.
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

#!/bin/bash

# Function definition
stop_services_with_name() {
    # Check if an argument was provided
    if [ -z "$1" ]; then
        echo "Usage: stop_services_with_name <full_path_to_process>"
        return 1
    fi

    # Use pgrep with -f option to find process IDs by full path
    # Note: macOS and most Linux distributions support these options
    local pids=$(pgrep -f "$1")

    # Check if any processes were found
    if [ -z "$pids" ]; then
        echo "No process found with the path: $1"
        return 0
    fi

    # Send the SIGTERM signal to each found process ID
    for pid in $pids; do
        echo "Sending SIGTERM to process ID $pid..."
        kill -15 "$pid"
        if [ $? -eq 0 ]; then
            echo "Process $pid has been terminated."
        else
            echo "Failed to terminate process $pid."
            return 1
        fi
    done

    return 0
}


#!/bin/bash

check_services_with_name() {
    local binary_path="$1"
  echo "ddddddddddd""$SUPPRESS_OUTPUT"
    pgrep -f "$binary_path" > /dev/null 2>&1


    if [ $? -eq 0 ]; then
        if [ -z "$SUPPRESS_OUTPUT" ]; then
            echo "A process with the path $binary_path is running."
        fi
        return 0
    else
        if [ -z "$SUPPRESS_OUTPUT" ]; then
            echo "No process found with the path $binary_path."
        fi
        return 1
    fi
}


# Example usage of check_services_with_name
# Replace "/full/path/to/binary" with the actual full path of the binary you want to check
# check_services_with_name "/full/path/to/binary"

