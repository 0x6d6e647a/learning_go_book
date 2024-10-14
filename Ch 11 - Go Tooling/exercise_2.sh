#!/usr/bin/env bash

declare -r staticcheck_url='honnef.co/go/tools/cmd/staticcheck@latest'

if ! type staticcheck &>/dev/null; then
    if ! go install "${staticcheck_url}"; then
        echo "go install staticcheck failed!"
        exit 1
    fi
fi

declare -r source='exercise_1.go'

staticcheck "${source}"
