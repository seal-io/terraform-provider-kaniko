terraform {
  required_providers {
    kaniko = {
      source = "registry.terraform.io/seal-io/kaniko"
    }
  }
}

provider "kaniko" {}

resource "kaniko_image" "example" {
  context     = "git://github.com/seal-io/simple-web-service"
  dockerfile  = "Dockerfile"
  destination = "docker.io/seal-io/test:1"

  build_arg = {
  }

  cache        = false
  no_push      = false
  reproducible = false
}
