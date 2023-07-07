#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
source "${ROOT_DIR}/hack/lib/init.sh"

PACKAGE_DIR="${ROOT_DIR}/.dist/package"

function release() {
  local prefix
  prefix=$(seal::target::build_prefix)

  local checksum_path="${PACKAGE_DIR}/${prefix}_${GIT_VERSION#v}_SHA256SUMS"
  shasum -a 256 "${PACKAGE_DIR}"/* | sed -e "s#${PACKAGE_DIR}/##g" >"${checksum_path}"
  if [[ -n "${GPG_FINGERPRINT:-}" ]]; then
    gpg --batch --local-user "${GPG_FINGERPRINT}" --detach-sign "${checksum_path}"
  else
    gpg --batch --detach-sign "${checksum_path}"
  fi
}

#
# main
#

seal::log::info "+++ PACKAGE +++" "tag: ${GIT_VERSION}"

release

seal::log::info "--- PACKAGE ---"
