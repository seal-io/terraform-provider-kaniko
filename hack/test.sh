#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${ROOT_DIR}/hack/lib/init.sh"

TEST_DIR="${ROOT_DIR}/.dist/test"
mkdir -p "${TEST_DIR}"

function test() {
  local tags=()
  # shellcheck disable=SC2086
  IFS=" " read -r -a tags <<<"$(seal::target::build_tags)"

  CGO_ENABLED=1 go test \
    -v \
    -failfast \
    -race \
    -cover \
    -timeout=30m \
    -tags="${tags[*]}" \
    -run="^Test[^(Acc)]" \
    -coverprofile="${TEST_DIR}/coverage.out" \
    "${ROOT_DIR}/..."
}

#
# main
#

seal::log::info "+++ TEST +++"

test

seal::log::info "--- TEST ---"
