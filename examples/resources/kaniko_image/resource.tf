resource "kaniko_image" "example" {
  context     = ""
  dockerfile  = ""
  destination = ""

  build_arg = {
  }

  cache        = false
  no_push      = false
  reproducible = false
}