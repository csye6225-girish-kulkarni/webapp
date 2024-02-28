packer {
  required_plugins {
    googlecompute = {
      source  = "github.com/hashicorp/googlecompute"
      version = "~> 1"
    }
  }
}

locals {
  timestamp = formatdate("YYYYMMDDhhmmss", timestamp())
}

source "googlecompute" "webapp-image" {
  project_id          = var.gcp_project_id
  source_image_family = "centos-stream-8"
  ssh_username        = "centos"
  image_name          = "webapp-image-${local.timestamp}"
  zone                = var.gcp_zone
  disk_size           = var.gcp_disk_size
  disk_type           = var.gcp_disk_type
}

build {
  sources = ["source.googlecompute.webapp-image"]

  provisioner "shell" {
    script = "./create_user.sh"
  }

  provisioner "shell" {
    script = "./install_go.sh"
  }

  provisioner "file" {
    source      = "webapp"
    destination = "/tmp/webapp"
  }

  provisioner "shell" {
    script = "./build_webapp.sh"
  }

  provisioner "file" {
    source      = "./webapp.service"
    destination = "/tmp/webapp.service"
  }

  provisioner "shell" {
    script = "./systemd_config.sh"
  }
}