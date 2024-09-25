# Terraform Provider for Nokia QIP

> **⚠️ Important Notice:** Vitesco Technologies has ceased using Nokia QIP as the main DNS tool and we currently don't use this project anymore. No further updates or maintenance will be provided for now.
>
> During our merging with Schaeffler, we will be re-evaluating this tool.
>
> Feel free to keep using the tool if you are able to support yourself. Add your comments to [#18](https://github.com/Vitesco-Technologies/terraform-provider-qip/issues/18).

![GitHub tag (with filter)](https://img.shields.io/github/v/release/Vitesco-Technologies/terraform-provider-qip)
[![Terraform Registry](https://img.shields.io/badge/Terraform_Registry-Vitesco--Technologies%2Fqip-blue)][terraform-registry]
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/Vitesco-Technologies/terraform-provider-qip/golang.yml)
![GitHub License](https://img.shields.io/github/license/Vitesco-Technologies/terraform-provider-qip)
[![Open Source](https://img.shields.io/badge/Vitesco_Technologies-open--source-yellow)](https://github.com/Vitesco-Technologies)

The provider for Nokia QIP will allow you to retrieve metadata from QIP or manage IPv4 addresses including their DNS names.

Documentation and releases can also be found on the [Terraform Registry under Vitesco-Technologies/qip][terraform-registry],
which also can be found as [Vitesco-Technologies/qip-address on the Terraform Registry](https://registry.terraform.io/modules/Vitesco-Technologies/qip-address/module/latest).

Features:

- Data sources for `qip_v4address` and `qip_v4subnet`
- Manage addresses with `qip_v4address`

Also see the Terraform module [qip-address](https://github.com/Vitesco-Technologies/terraform-module-qip-address).

Build based on the Swagger API documentation that should be available with your QIP instance: `https://qip.example.com.com/rest-api/`

## How to use

Please see the [documentation on the Terraform registry][terraform-registry] for some examples.

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

[terraform-registry]: https://registry.terraform.io/providers/Vitesco-Technologies/qip/latest
