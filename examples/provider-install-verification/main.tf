terraform {
  required_providers {
    kaniko = {
      source = "registry.terraform.io/seal-io/kaniko"
    }
  }
}

provider "kaniko" {}

