#!/bin/bash

# CodeKiwi Configuration Loader
# This script loads the config.env file and provides helper functions

# Determine the project root directory
get_project_root() {
    local script_path="${BASH_SOURCE[0]}"
    local script_dir="$(cd "$(dirname "$script_path")" && pwd)"
    echo "$(cd "$script_dir/.." && pwd)"
}

# Load configuration
load_config() {
    local project_root="$(get_project_root)"
    local config_file="$project_root/config.env"

    if [ ! -f "$config_file" ]; then
        echo "Error: config.env not found at $config_file" >&2
        return 1
    fi

    # Source the config file
    set -a  # automatically export all variables
    source "$config_file"
    set +a

    # Compute derived values
    export CODEKIWI_INSTALL_DIR="$HOME/$CODEKIWI_INSTALL_DIR_NAME"
    export CODEKIWI_FULL_IMAGE_NAME="$CODEKIWI_IMAGE_REGISTRY/$CODEKIWI_IMAGE_NAME"
}

# Get full image name with tag
get_image_with_tag() {
    local tag="${1:-$CODEKIWI_IMAGE_TAG_DEFAULT}"
    echo "$CODEKIWI_FULL_IMAGE_NAME:$tag"
}

# Get install directory (supports both dev and production mode)
get_install_dir() {
    local script_dir="$1"
    local project_root="$(cd "$script_dir/../.." 2>/dev/null && pwd)"

    if [ -f "$project_root/$CODEKIWI_COMPOSE_FILE_DEV" ]; then
        # Development mode
        echo "$project_root"
    else
        # Production mode
        echo "$CODEKIWI_INSTALL_DIR"
    fi
}

# Get instances directory
get_instances_dir() {
    local install_dir="$1"
    echo "$install_dir/instances"
}

# Load configuration on source
load_config
