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
OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

source $OPENIM_ROOT/scripts/style_info.sh
source $OPENIM_ROOT/scripts/path_info.sh

cd "$SCRIPTS_ROOT"

service_port_name=(
 openImChatApiPort
 openImAdminApiPort
   #api port name
   openImAdminPort
   openImChatPort
)

for i in ${service_port_name[*]}; do
  list=$(cat $config_path | grep -w ${i} | awk -F '[:]' '{print $NF}')
  list_to_string $list
  for j in ${ports_array}; do
    echo -e "========${j}======="
      name="ps -aux |grep -w $i |grep -v grep"
      count="${name}| wc -l"
      echo -e "==========$(eval ${count})=========="
      if [ $(eval ${count}) -gt 0 ]; then
        pid="${name}| awk '{print \$2}'"
        echo -e "${SKY_BLUE_PREFIX}Killing service:$i pid:$(eval $pid)${COLOR_SUFFIX}"
        #kill the service that existed
        kill -9 $(eval $pid)
        echo -e "${SKY_BLUE_PREFIX}service:$i was killed ${COLOR_SUFFIX}"
      fi
  done
done