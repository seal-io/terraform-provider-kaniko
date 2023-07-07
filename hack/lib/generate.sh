#!/usr/bin/env bash

# -----------------------------------------------------------------------------
# Lint variables helpers. These functions need the
# following variables:
#
#    TFPLUGINDOCS_VERSION  -  The Terraform docs plugin version, default is v1.50.1.

tfplugindocs_version=${TFPLUGINDOCS_VERSION:-"v0.14.1"}

function seal::docs::tfplugin::install() {
  GOBIN="${ROOT_DIR}/.sbin" go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@"${tfplugindocs_version}"
}

function seal::docs::tfplugin::validate() {
  # shellcheck disable=SC2046
  if [[ -n "$(command -v $(seal::docs::tfplugin::bin))" ]]; then
    return 0
  fi

  seal::log::info "installing tfplugindocs"
  if seal::docs::tfplugin::install; then
    return 0
  fi
  seal::log::error "no tfplugindocs available"
  return 1
}

function seal::docs::tfplugin::bin() {
  local bin="tfplugindocs"
  if [[ -f "${ROOT_DIR}/.sbin/tfplugindocs" ]]; then
    bin="${ROOT_DIR}/.sbin/tfplugindocs"
  fi
  echo -n "${bin}"
}

function seal::docs::generate() {
  if ! seal::docs::tfplugin::validate; then
    seal::log::fatal "cannot execute tfplugindocs as client is not found"
  fi

  seal::log::debug "tfplugindocs generate $*"
  $(seal::docs::tfplugin::bin) generate "$@"
}
