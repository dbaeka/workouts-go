#!/bin/bash

readonly env_file="$2"
readonly service="$1"

cd "./internal/$service" || exit
env "$(cat "../../.env" "../../$env_file" | grep -Ev '^#' | xargs)" go test -count=1 -race ./...
