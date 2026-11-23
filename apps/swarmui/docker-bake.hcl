target "docker-metadata-action" {}

variable "APP" {
  default = "swarmui"
}

variable "VERSION" {
  // Pinned manually to upstream master commit 66ac87f (2026-11-22).
  // Why not a Renovate-tracked github-release: SwarmUI has published nothing
  // newer than 0.9.8-Beta (2026-02-06), whose Model Downloader cannot resolve
  // civitai.com/api/download/models/<id>?fileId=<id> links — it 404s the wrong
  // metadata endpoint (/api/v1/models/<versionId>?fileId=...). The fix exists
  // only on master so far. Renovate tracking is intentionally disabled here;
  // restore the comment
  // "# renovate: datasource=github-releases depName=mcmonkeyprojects/SwarmUI"
  // directly above the default once upstream tags a release that includes the
  // rewritten downloader.
  default = "66ac87fcc41b5163d6b372658bdf85374a4c8158"
}

variable "SOURCE" {
  default = "https://github.com/mcmonkeyprojects/SwarmUI"
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
