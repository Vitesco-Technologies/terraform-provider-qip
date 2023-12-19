provider "qip" {
  server   = "https://qip.example.com"
  org      = "Example"
  username = "admin"
  password = "admin"
}

terraform {
  required_providers {
    qip = {
      source  = "Vitesco-Technologies/qip"
      version = ">0"
    }
  }
}
