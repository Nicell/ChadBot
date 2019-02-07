#!/bin/bash

if [ $# -ne 2 ]; then
    echo "Missing required two arguments (DISCORD, YOUTUBE), given $# arguments"
    exit 1
fi

declare -r IMAGE_NAME="nicell/chadbot"
declare -r IMAGE_TAG="latest"

echo "Starting container for image '$IMAGE_NAME:$IMAGE_TAG'"
docker run -e DISCORD="$1" -e YOUTUBE="$2" $IMAGE_NAME:$IMAGE_TAG
