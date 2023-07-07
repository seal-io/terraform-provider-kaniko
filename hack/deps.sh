#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${ROOT_DIR}/hack/lib/init.sh"

function mod() {
  go mod tidy
  go mod download
}

#
# main
#

seal::log::info "+++ MOD +++"

mod

seal::log::info "--- MOD ---"
