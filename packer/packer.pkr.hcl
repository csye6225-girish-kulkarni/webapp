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
  machine_type        = "n1-standard-2"
}

build {
  sources = ["source.googlecompute.webapp-image"]

  provisioner "shell" {
    script = "./scripts/create_user.sh"
  }

  provisioner "shell" {
    script = "./scripts/install_go.sh"
  }

  provisioner "file" {
    source      = "config.yaml"
    destination = "/tmp/config.yaml"
  }

  provisioner "shell" {
    script = "./scripts/install_ops_agent.sh"
  }

  provisioner "file" {
    source      = "webapp"
    destination = "/tmp/webapp"
  }

  provisioner "shell" {
    script = "./scripts/build_webapp.sh"
  }

  provisioner "file" {
    source      = "./webapp.service"
    destination = "/tmp/webapp.service"
  }

  provisioner "shell" {
    script = "./scripts/systemd_config.sh"
  }
}