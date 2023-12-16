#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

# Include shell font styles and some basic information
SCRIPTS_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

config_file="${OPENIM_ROOT}/config/config.yaml"

# Initialize flags
FORCE=false
SKIP=false

show_help() {
    echo "Usage: init-config.sh [options]"
    echo "Options:"
    echo "  -h, --help             Show this help message"
    echo "  --force                Overwrite existing files without prompt"
    echo "  --skip                 Skip generation if file exists"
    echo "  --clean-config         Clean all configuration files"
}

clean_config() {
    echo "Cleaning configuration files..."
    rm -f "${config_file}"
    echo "Configuration files cleaned."
}

generate_config() {
    echo "Generating configuration file..."
    cp "${OPENIM_ROOT}/deployments/templates/config.yaml" "${config_file}"
    echo "Configuration file generated."
}

overwrite_prompt() {
    while true; do
        read -p "Configuration file exists. Overwrite? [Y/N]: " yn
        case $yn in
            [Yy]* ) generate_config; break;;
            [Nn]* ) echo "Skipping generation."; exit;;
            * ) echo "Please answer yes or no.";;
        esac
    done
}

# Parse command line arguments
for i in "$@"
do
case $i in
    -h|--help)
    show_help
    exit 0
    ;;
    --force)
    FORCE=true
    shift
    ;;
    --skip)
    SKIP=true
    shift
    ;;
    --clean-config)
    clean_config
    exit 0
    ;;
    *)
    # unknown option
    show_help
    exit 1
    ;;
esac
done

if [[ "${FORCE}" == "true" ]]; then
    generate_config
elif [[ "${SKIP}" == "true" ]] && [[ -f "${config_file}" ]]; then
    echo "Configuration file already exists. Skipping generation."
else
    if [[ -f "${config_file}" ]]; then
        overwrite_prompt
    else
        generate_config
    fi
fi
