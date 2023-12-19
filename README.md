# Terraform Provider for Nokia QIP

The provider for Nokia QIP will allow you to retrieve metadata from QIP or manage IPv4 addresses including their DNS names.

Features:

- Data sources for `qip_v4address` and `qip_v4subnet`
- Manage addresses with `qip_v4address`

<!-- TODO: when published
Also see the Terraform module [qip-address](https://github.com/Vitesco-Technologies/terraform-module-qip-address).
-->

Build based on the Swagger API documentation that should be available with your QIP instance: `https://qip.example.com.com/rest-api/`

## How to use

Please see the [documentation](docs/) for some examples.

Very basic usage:

```terraform
resource "qip_v4address" "address" {
  address = "192.0.2.23"
  subnet  = "192.0.2.0"

  name         = "my-example"
  // object_class = "Virtual Server"
  // description  = "Example System"
  // domain_name  = "corp.example.com"
}

resource "qip_v4address" "address" {
  // selecting a free address in the subnet
  subnet = "192.0.2.0"

  subnet_range_start = "192.0.2.30"
  subnet_range_end = "192.0.2.50"

  name   = "my-example"
}

data "qip_v4subnet" "subnet" {
  address = "192.0.2.0"
}

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
```

## Installation for testing purposes

First, build and install the provider locally.

```shell
make install
```

Then, run the following command to initialize the workspace and apply the sample configuration.

```shell
terraform init && terraform apply
```
