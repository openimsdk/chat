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
    exec >> ${SCRIPTS_ROOT}/../logs/openIM.log 2>&1
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
  grep -oP "${key}: \[\s*\K[^\]]+" "$config_yaml" | xargs
}

# Extract port numbers from the YAML configuration
openImChatApiPort=$(extract_yaml_value 'openImChatApiPort')
openImAdminApiPort=$(extract_yaml_value 'openImAdminApiPort')
openImAdminPort=$(extract_yaml_value 'openImAdminPort')
openImChatPort=$(extract_yaml_value 'openImChatPort')

for i in "${service_port_name[@]}"; do
  case $i in
    "openImChatApiPort")
      new_service_name="chat-api"
      new_service_port=$openImChatApiPort
      ;;
    "openImAdminApiPort")
      new_service_name="admin-api"
      new_service_port=$openImAdminApiPort
      ;;
    "openImAdminPort")
      new_service_name="admin-rpc"
      new_service_port=$openImAdminPort
      ;;
    "openImChatPort")
      new_service_name="chat-rpc"
      new_service_port=$openImChatPort
      ;;
    *)
      echo "Invalid service name: $i"
      exit -1
      ;;
  esac


  ports=$(ss -tunlp | grep "$new_service_name" | awk '{print $5}' | awk -F '[:]' '{print $NF}')

found_port=false
for port in $ports; do
  if [[ "$port" == "$new_service_port" ]]; then
    echo -e "${new_service_port}${GREEN_PREFIX} port has been listening, belongs service is ${i}${COLOR_SUFFIX}"
    found_port=true
    break
  fi
done

if [[ "$found_port" != true ]]; then
  echo -e "${YELLOW_PREFIX}${i}${COLOR_SUFFIX}${RED_PREFIX} service does not start normally, expected port is ${COLOR_SUFFIX}${YELLOW_PREFIX}${new_service_port}${COLOR_SUFFIX}"
  echo -e "${RED_PREFIX}please check ${SCRIPTS_ROOT}/../logs/openIM.log ${COLOR_SUFFIX}"
  exit -1
fi

done

