#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${ROOT_DIR}/hack/lib/init.sh"

BUILD_DIR="${ROOT_DIR}/.dist/build"
rm -rf "${BUILD_DIR}"
mkdir -p "${BUILD_DIR}"

function build() {
  local prefix
  prefix=$(seal::target::build_prefix)

  local ldflags=(
    "-X main.version=${GIT_VERSION}"
    "-X main.commit=${GIT_COMMIT}"
    "-w -s"
    "-extldflags '-static'"
  )

  local tags=()
  # shellcheck disable=SC2086
  IFS=" " read -r -a tags <<<"$(seal::target::build_tags)"

  local platforms=()
  # shellcheck disable=SC2086
  IFS=" " read -r -a platforms <<<"$(seal::target::build_platforms)"

  for platform in "${platforms[@]}"; do
    local os_arch
    IFS="/" read -r -a os_arch <<<"${platform}"

    local os="${os_arch[0]}"
    local arch="${os_arch[1]}"
    local ext=""
    if [[ "${platform}" =~ windows/* ]]; then
      ext=".exe"
    fi

    GOOS=${os} GOARCH=${arch} CGO_ENABLED=0 go build \
      -trimpath \
      -ldflags="${ldflags[*]}" \
      -tags="${os} ${tags[*]}" \
      -o="${BUILD_DIR}/${os}/${arch}/${prefix}_${GIT_VERSION}${ext}" \
      "${ROOT_DIR}"
  done
}

#
# main
#

seal::log::info "+++ BUILD +++" "info: ${GIT_VERSION},${GIT_COMMIT:0:7},${GIT_TREE_STATE},${BUILD_DATE}"

build

seal::log::info "--- BUILD ---"
