---
layout: ""
page_title: "Kaniko Provider"
description: The Kaniko provider for Terraform is a plugin that leverage GoogleContainerTools/Kaniko to build container images from a Dockerfile, inside a container or Kubernetes cluster.
---

# Kaniko Provider

The Kaniko provider for Terraform is a plugin that leverage [GoogleContainerTools/Kaniko](https://github.com/GoogleContainerTools/kaniko) to build container images from a Dockerfile, inside a container or Kubernetes cluster.

## Example Usage

```terraform
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
```
