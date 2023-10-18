provider "qip" {
  server   = "https://qip.example.com"
  org      = "Example"
  username = "admin"
  password = "admin"
}

terraform {
  required_providers {
    qip = {
      source  = "registry.terraform.io/vitesco-technologies/qip"
      version = ">0"
    }
  }
}
