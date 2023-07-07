# Kaniko for Terraform [![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/seal-io/terraform-provider-kaniko?label=release)](https://github.com/seal-io/terraform-provider-kaniko/releases) [![license](https://img.shields.io/github/license/seal-io/terraform-provider-kaniko.svg)]()

The Kaniko provider for Terraform is a plugin that leverage [GoogleContainerTools/Kaniko](https://github.com/GoogleContainerTools/kaniko) to build container images from a Dockerfile, inside a container or Kubernetes cluster.

This provider is maintained by [Seal](https://github.com/seal-io).

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) 1.4.x
-	[Go](https://golang.org/doc/install) 1.19.x (to build the provider plugin)
