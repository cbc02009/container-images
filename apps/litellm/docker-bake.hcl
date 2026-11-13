target "docker-metadata-action" {}

variable "APP" {
  default = "litellm"
}

variable "VERSION" {
  // renovate: datasource=docker depName=ghcr.io/berriai/litellm
  default = "v1.109.0-rc.1"
}

variable "SOURCE" {
  default = "https://github.com/BerriAI/litellm"
}

group "default" {
  targets = ["image-local"]
}

target "image" {
  inherits = ["docker-metadata-action"]
  args = {
    VERSION = "${VERSION}"
  }
  labels = {
    "org.opencontainers.image.source" = "${SOURCE}"
  }
}

target "image-local" {
  inherits = ["image"]
  output = ["type=docker"]
  tags = ["${APP}:${VERSION}"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/amd64"
  ]
}
