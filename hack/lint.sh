#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${ROOT_DIR}/hack/lib/init.sh"

function check_dirty() {
  [[ "${LINT_DIRTY:-false}" == "true" ]] || return 0

  if [[ -n "$(command -v git)" ]]; then
    if git_status=$(git status --porcelain 2>/dev/null) && [[ -n ${git_status} ]]; then
      seal::log::fatal "the git tree is dirty:\n$(git status --porcelain)"
    fi
  fi
}

function lint() {
  local path="$1"
  shift 1

  seal::format::run "${path}"
  GOLANGCI_LINT_CACHE="$(go env GOCACHE)" seal::lint::run --build-tags="$*" "${path}/..."
}

function after() {
  check_dirty
}

#
# main
#

seal::log::info "+++ LINT +++"

lint "${ROOT_DIR}" "$(seal::target::build_tags)"

after

seal::log::info "--- LINT ---"
