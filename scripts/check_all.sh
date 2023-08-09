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

#Include shell font styles and some basic information
SCRIPTS_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
OPENIM_ROOT=$(dirname "${SCRIPTS_ROOT}")/..

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

switch=$(cat $config_path | grep demoswitch | awk -F '[:]' '{print $NF}')
for i in ${service_port_name[*]}; do
  case $i in
    "openImChatApiPort")
      new_service_name="chat-api"
      new_service_port="10008"
      ;;
    "openImAdminApiPort")
      new_service_name="admin-rpc"
      new_service_port="30200"
      ;;
    "openImAdminPort")
      new_service_name="chat-rpc"
      new_service_port="30300"
      ;;
    "openImChatPort")
      new_service_name="admin-api"
      new_service_port="10009"
      ;;
    *)
      echo "Invalid service name: $i"
      exit -1
      ;;
  esac

  port=$(ss -tunlp | grep "$new_service_name" | awk '{print $5}' | awk -F '[:]' '{print $NF}')
  if [[ "$port" != "$new_service_port" ]]; then
    echo -e "${YELLOW_PREFIX}${i}${COLOR_SUFFIX}${RED_PREFIX} service does not start normally, not initiated port is ${COLOR_SUFFIX}${YELLOW_PREFIX}${new_service_port}${COLOR_SUFFIX}"
    echo -e "${RED_PREFIX}please check ${SCRIPTS_ROOT}/../logs/openIM.log ${COLOR_SUFFIX}"
    exit -1
  else
    echo -e "${new_service_port}${GREEN_PREFIX} port has been listening, belongs service is ${i}${COLOR_SUFFIX}"
  fi
done
