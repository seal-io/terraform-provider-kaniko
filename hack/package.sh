#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${ROOT_DIR}/hack/lib/init.sh"

BUILD_DIR="${ROOT_DIR}/.dist/build"
PACKAGE_DIR="${ROOT_DIR}/.dist/package"
rm -rf "${PACKAGE_DIR}"
mkdir -p "${PACKAGE_DIR}"

function package() {
  local prefix
  prefix=$(seal::target::build_prefix)

  local platforms=()
  # shellcheck disable=SC2086
  IFS=" " read -r -a platforms <<<"$(seal::target::build_platforms)"

  # archive
  for platform in "${platforms[@]}"; do
    local os_arch
    IFS="/" read -r -a os_arch <<<"${platform}"

    local os="${os_arch[0]}"
    local arch="${os_arch[1]}"
    local ext=""
    if [[ "${platform}" =~ windows/* ]]; then
      ext=".exe"
    fi

    local src="${BUILD_DIR}/${os}/${arch}/${prefix}_${GIT_VERSION}${ext}"
    local dst="${PACKAGE_DIR}/${prefix}_${GIT_VERSION#v}_${os}_${arch}.zip"
    zip -1qj "${dst}" "${src}"
  done

  # manifest
  cp -f "${ROOT_DIR}/terraform-registry-manifest.json" "${PACKAGE_DIR}/${prefix}_${GIT_VERSION#v}_manifest.json"
}

#
# main
#

seal::log::info "+++ PACKAGE +++" "tag: ${GIT_VERSION}"

package

seal::log::info "--- PACKAGE ---"
