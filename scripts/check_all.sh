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
echo "# Start Chat check_all.sh ${DATA}, For local deployments, use ./check_all.sh --print-screen"

# 检查第一个参数是否为 --print-screen
if [ "$1" == "--print-screen" ]; then
    PRINT_SCREEN=1
fi

mkdir -p ${SCRIPTS_ROOT}/../logs
# 如果没有设置 PRINT_SCREEN 标记，那么进行日志重定向
if [ -z "$PRINT_SCREEN" ]; then
    exec >> ${SCRIPTS_ROOT}/../logs/openim_$(date '+%Y%m%d').log 2>&1
fi

#Include shell font styles and some basic information
source $SCRIPTS_ROOT/style_info.sh
source $SCRIPTS_ROOT/path_info.sh
source $SCRIPTS_ROOT/function.sh

service_port_name=(
 openImChatApiPort
 openImAdminApiPort
   #api port name
   openImAdminPort
   openImChatPort
)
sleep 10


# Define the path to the configuration YAML file
config_yaml="$OPENIM_ROOT/config/config.yaml" # Replace with the actual path to your YAML file

# Function to extract a value from the YAML file and remove any leading/trailing whitespace
extract_yaml_value() {
  local key=$1

  # Detect the operating system
  case "$(uname)" in
    "Linux")
      # Use grep with Perl-compatible regex for Linux
      grep -oP "${key}: \[\s*\K[^\]]+" "$config_yaml" | xargs
      ;;
    "Darwin")
      # Use sed for macOS
      sed -nE "/${key}: \[ */{s///; s/\].*//; p;}" "$config_yaml" | tr -d '[]' | xargs
      ;;
    *)
      echo "Unsupported operating system"
      exit 1
      ;;
  esac
}

# Extract port numbers from the YAML configuration
declare -A service_ports=(
  ["openImChatApiPort"]="chat-api"
  ["openImAdminApiPort"]="admin-api"
  ["openImAdminPort"]="admin-rpc"
  ["openImChatPort"]="chat-rpc"
)

for i in "${!service_ports[@]}"; do
  service_port=$(extract_yaml_value "$i")
  new_service_name=${service_ports[$i]}

  # Check for empty port value
  if [ -z "$service_port" ]; then
    echo "No port value found for $i"
    continue
  fi

  # Determine command based on OS
  case "$(uname)" in
    "Linux")
      ports=$(ss -tunlp | grep "$new_service_name" | awk '{print $5}' | awk -F '[:]' '{print $NF}')
      ;;
    "Darwin")
      ports=$(lsof -i -P | grep LISTEN | grep "$new_service_name" | awk '{print $9}' | awk -F '[:]' '{print $2}')
      ;;
    *)
      echo "Unsupported operating system"
      exit 1
      ;;
  esac

  found_port=false
  for port in $ports; do
    if [[ "$port" == "$service_port" ]]; then
      echo -e "${service_port}${GREEN_PREFIX} port has been listening, belongs service is ${new_service_name}${COLOR_SUFFIX}"
      found_port=true
      break
    fi
  done

  if [[ "$found_port" != true ]]; then
    echo -e "${YELLOW_PREFIX}${new_service_name}${COLOR_SUFFIX}${RED_PREFIX} service does not start normally, expected port is ${COLOR_SUFFIX}${YELLOW_PREFIX}${service_port}${COLOR_SUFFIX}"
    echo -e "${RED_PREFIX}please check ${SCRIPTS_ROOT}/../logs/chat_$(date '+%Y%m%d').log ${COLOR_SUFFIX}"
    exit -1
  fi
done
