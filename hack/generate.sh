#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${ROOT_DIR}/hack/lib/init.sh"

function generate() {
  seal::docs::generate
}

#
# main
#

seal::log::info "+++ GENERATE +++"

generate

seal::log::info "--- GENERATE ---"
