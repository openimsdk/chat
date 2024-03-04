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

source $SCRIPTS_ROOT/style-info.sh
source $SCRIPTS_ROOT/path-info.sh
source $SCRIPTS_ROOT/function.sh
source $SCRIPTS_ROOT/util.sh

# if [ ! -d "${OPENIM_ROOT}/_output/bin/platforms" ]; then
#   cd $OPENIM_ROOT
#   # exec build-all-service.sh
#   "${SCRIPTS_ROOT}/build-all-service.sh"
# fi

bin_dir="$BIN_DIR"
logs_dir="$SCRIPTS_ROOT/../_output/logs"

# Define the path to the configuration file
CONFIG_FILE="${OPENIM_ROOT}/config/config.yaml"

# Check if the configuration file exists
if [ -f "$CONFIG_FILE" ]; then
    # The file exists
    echo "Configuration file already exists at $CONFIG_FILE."
else
    echo ""
    # The file does not exist
    echo "Error: Configuration file does not exist."
    echo "+++ You need to execute 'make init' to generate the configuration file and then modify the configuration items."
    echo ""
    exit 1
fi

#service filename
service_filename=(
  chat-api
  admin-api
  #rpc
  admin-rpc
  chat-rpc
)

#service config port name
service_port_name=(
openImChatApiPort
openImAdminApiPort
  #api port name
  openImAdminPort
  openImChatPort
)

service_prometheus_port_name=(

)

# Automatically created when there is no bin, logs folder
if [ ! -d $logs_dir ]; then
  mkdir -p $logs_dir
fi
cd $SCRIPTS_ROOT

for ((i = 0; i < ${#service_filename[*]}; i++)); do
  rm -rf ${logs_dir}/chat_tmp_$(date '+%Y%m%d').log
  #Check whether the service exists
#  service_name="ps |grep -w ${service_filename[$i]} |grep -v grep"
#  count="${service_name}| wc -l"
#
#  if [ $(eval ${count}) -gt 0 ]; then
#    pid="${service_name}| awk '{print \$2}'"
#    echo  "${service_filename[$i]} service has been started,pid:$(eval $pid)"
#    echo  "killing the service ${service_filename[$i]} pid:$(eval $pid)"
#    #kill the service that existed
#    kill -9 $(eval $pid)
#    sleep 0.5
#  fi
  cd $SCRIPTS_ROOT

  #Get the rpc port in the configuration file
  portList=$(cat $config_path | grep ${service_port_name[$i]} | awk -F '[:]' '{print $NF}')
  list_to_string ${portList}
  service_ports=($ports_array)


  #Start related rpc services based on the number of ports
  for ((j = 0; j < ${#service_ports[*]}; j++)); do
    if [ ! -e "$bin_dir/${service_filename[$i]}" ]; then
      echo -e  ${RED_PREFIX}"Error: ${service_filename[$i]} does not exist,Start fail!"${COLOR_SUFFIX}
      echo "start build these binary"
      "./build-all-service.sh"
    fi
    #Start the service in the background
    cmd="$bin_dir/${service_filename[$i]} -port ${service_ports[$j]} --config_folder_path ${config_path}"
    if [ $i -eq 0 -o $i -eq 1 ]; then
      cmd="$bin_dir/${service_filename[$i]} -port ${service_ports[$j]} --config_folder_path ${config_path}"
    fi
    echo $cmd
    nohup $cmd >> ${logs_dir}/chat_$(date '+%Y%m%d').log 2> >(tee -a ${logs_dir}/chat_err_$(date '+%Y%m%d').log ${logs_dir}/chat_tmp_$(date '+%Y%m%d').log) &
  done
done


sleep 1

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





