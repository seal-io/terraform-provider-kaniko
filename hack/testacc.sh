#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${ROOT_DIR}/hack/lib/init.sh"

TEST_DIR="${ROOT_DIR}/.dist/test"
mkdir -p "${TEST_DIR}"

function testacc() {
  TF_ACC=1 go test \
    -v \
    -failfast \
    -run="^TestAcc[A-Z]+" \
    -timeout=30m \
    "${ROOT_DIR}/..."
}

#
# main
#

seal::log::info "+++ TEST ACC +++"

testacc

seal::log::info "--- TEST ACC ---"
