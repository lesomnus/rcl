#!/usr/bin/env bash

set -o errexit
set -o pipefail
# set -o xtrace

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)" # Directory where this script exists.
__root="$(cd "$(dirname "${__dir}")" && pwd)"         # Root directory of project.

cd "${__root}"

mkdir -p ".pc"
go run ./cmd/gen-pc -o ".pc/rcl-go.pc"
