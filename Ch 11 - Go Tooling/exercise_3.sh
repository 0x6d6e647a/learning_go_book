#!/usr/bin/env bash

declare -r target_os='windows'
declare -r target_arch='amd64'
declare -r source='exercise_1.go'

GOOS="${target_os}" GOARCH="${target_arch}" go build "${source}"
