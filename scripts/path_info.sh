#!/usr/bin/env bash

#Don't put the space between "="

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
SCRIPTS_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

demo_server_name="chat-api"
demo_server_binary_root="${OPENIM_ROOT}/bin/"

# Determine the architecture and version
architecture=$(uname -m)
version=$(uname -s | tr '[:upper:]' '[:lower:]')

#Include shell font styles and some basic information
source $SCRIPTS_ROOT/style_info.sh

cd $SCRIPTS_ROOT

# Define the supported architectures and corresponding bin directories
declare -A supported_architectures=(
    ["linux-amd64"]="_output/bin/platforms/linux/amd64"
    ["linux-arm64"]="_output/bin/platforms/linux/arm64"
    ["linux-mips64"]="_output/bin/platforms/linux/mips64"
    ["linux-mips64le"]="_output/bin/platforms/linux/mips64le"
    ["linux-ppc64le"]="_output/bin/platforms/linux/ppc64le"
    ["linux-s390x"]="_output/bin/platforms/linux/s390x"
    ["darwin-amd64"]="_output/bin/platforms/darwin/amd64"
    ["windows-amd64"]="_output/bin/platforms/windows/amd64"
    ["linux-x86_64"]="_output/bin/platforms/linux/amd64"  # Alias for linux-amd64
    ["darwin-x86_64"]="_output/bin/platforms/darwin/amd64"  # Alias for darwin-amd64
)

# Check if the architecture and version are supported
if [[ -z ${supported_architectures["$version-$architecture"]} ]]; then
    echo -e "${BLUE_PREFIX}================> Unsupported architecture: $architecture or version: $version${COLOR_SUFFIX}"
    exit 1
fi

echo -e "${BLUE_PREFIX}================> Architecture: $architecture${COLOR_SUFFIX}"

# Set the BIN_DIR based on the architecture and version
BIN_DIR=${supported_architectures["$version-$architecture"]}

echo -e "${BLUE_PREFIX}================> BIN_DIR: $OPENIM_ROOT/$BIN_DIR${COLOR_SUFFIX}"

#Global configuration file default dir
config_path="$SCRIPTS_ROOT/../config/config.yaml"
configfile_path="$OPENIM_ROOT/config"
log_path="$SCRIPTS_ROOT/../log"

#servicefile dir path
service_source_root=(
  #api service file
  $OPENIM_ROOT/cmd/api/chat-api/
  $OPENIM_ROOT/cmd/api/admin-api/
  #rpc service file
  $OPENIM_ROOT/cmd/rpc/admin-rpc/
  $OPENIM_ROOT/cmd/rpc/chat-rpc/
)
#service filename
service_names=(
  #api service filename
  chat-api
  admin-api
  #rpc service filename
  admin-rpc
  chat-rpc
)
