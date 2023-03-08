terraform {
  required_providers {
    kaniko = {
      source = "example.com/gitlawr/kaniko"
    }
  }
}

provider "kaniko" {}

resource "kaniko_image" "example" {
  context = "git://github.com/gitlawr/go-example"
  dockerfile = "Dockerfile"
  destination = "docker.io/lawr/test:1"

  build_arg  = {
  }

  cache = false
  no_push = false
  reproducible = false
}
