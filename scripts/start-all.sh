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

export SUPPRESS_OUTPUT=1
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
  openImAdminPort
  openImChatPort
)

service_prometheus_port_name=(

)


#!/bin/bash

# Reusing stop_services_with_name and check_services_with_name functions as provided

check_and_stop_services() {
    local services=("$@")
    local service_stopped=0
    local attempts=0

    # Step 1: Check and stop each service if running
    for service in "${services[@]}"; do
        stop_services_with_name "$service" >/dev/null 2>&1
        if [ $? -eq 0 ]; then
            echo "Service running: $service. Attempting to stop."
            stop_services_with_name "$service"
        fi
    done


    # Step 2: Verify all services are stopped, retry up to 15 times if necessary
    while [ $attempts -lt 15 ]; do
        service_stopped=1

        for service in "${services[@]}"; do
            result=$(check_services_with_name "$service")
            if [ $? -eq 0 ]; then
                service_stopped=0
                break
            fi
        done
        if [ $service_stopped -eq 1 ]; then
            echo "All services have been successfully stopped."
            return 0
        fi

        sleep 1
        ((attempts++))
    done

    if [ $service_stopped -eq 0 ]; then
        echo "Failed to stop all services after 15 seconds."
        return 1
    fi
}



# Call the function with your full binary paths
check_and_stop_services "${binary_full_paths[@]}"
exit_status=$?

# Check the exit status and proceed accordingly
if [ $exit_status -eq 0 ]; then
    echo "Execution can continue."
else
    echo "Exiting due to failure in stopping services."
    exit 1
fi







# Automatically created when there is no bin, logs folder
if [ ! -d $logs_dir ]; then
  mkdir -p $logs_dir
fi
cd $SCRIPTS_ROOT

rm -rf ${logs_dir}/chat_tmp_$(date '+%Y%m%d').log
LOG_FILE=${logs_dir}/chat_$(date '+%Y%m%d').log
STDERR_LOG_FILE=${logs_dir}/chat_err_$(date '+%Y%m%d').log
TMP_LOG_FILE=${logs_dir}/chat_tmp_$(date '+%Y%m%d').log

cmd="${component_binary_full_paths} --config_folder_path ${config_path}"
echo $cmd ...............
nohup ${cmd} >> "${LOG_FILE}" 2> >(tee -a "${STDERR_LOG_FILE}" "$TMP_LOG_FILE" | while read line; do echo -e "\e[31m${line}\e[0m"; done >&2)
if [ $? -eq 0 ]; then
    echo "All components checked successfully"
    # Add the commands that should be executed next if the binary component was successful
else
    echo "Component check failed, program exiting"
    exit 1
fi

for ((i = 0; i < ${#service_filename[*]}; i++)); do

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


    nohup ${cmd} >> "${LOG_FILE}" 2> >(tee -a "${STDERR_LOG_FILE}" "$TMP_LOG_FILE" | while read line; do echo -e "\e[31m${line}\e[0m"; done >&2) &


  done
done




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
    echo -e "\033[0;32mAll chat services startup successful\033[0m"
fi

all_ports_listening=true


ports=(
  $(sed -n 's/.*openImChatApiPort: \[\(.*\)\].*/\1/p' ${config_path}/config.yaml)
  $(sed -n 's/.*openImAdminApiPort: \[\(.*\)\].*/\1/p' ${config_path}/config.yaml)
  $(sed -n 's/.*openImAdminPort: \[\(.*\)\].*/\1/p' ${config_path}/config.yaml)
  $(sed -n 's/.*openImChatPort: \[\(.*\)\].*/\1/p' ${config_path}/config.yaml)
)




for port in "${ports[@]}"; do
  if ! check_services_with_port "$port"; then
    all_ports_listening=false
    break
  fi
done

if $all_ports_listening; then
  echo "successful"
else
  echo "failed"
fi






