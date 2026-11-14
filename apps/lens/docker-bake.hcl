target "docker-metadata-action" {}

variable "APP" {
  default = "lens"
}

variable "VERSION" {
  // renovate: datasource=github-tags depName=BerriAI/lens
  default = "v1.108.0-rc.1"
}

variable "SOURCE" {
  default = "https://github.com/BerriAI/lens"
}

group "default" {
  targets = ["image-local"]
}

target "image" {
  inherits = ["docker-metadata-action"]
  // The Dockerfile's final stage; skips the upstream smoke-test stage, which
  // is exercised by container_test.go in PRs instead.
  target = "runtime"
  args = {
    LENS_REF     = "${VERSION}"
    LENS_VERSION = "${VERSION}"
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
