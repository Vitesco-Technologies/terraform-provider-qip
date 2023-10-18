# We use an address to derive the current domain and address from.
resource "qip_v4address" "address" {
  subnet = "192.0.2.0"
  name   = "my-example"
}

resource "qip_v4address_rr" "wildcard" {
  name        = "*.my-example"
  address     = qip_v4address.address.address
  domain_name = qip_v4address.address.domain_name
}

resource "qip_v4address_rr" "other_record" {
  name        = "other-example"
  address     = qip_v4address.address.address
  domain_name = qip_v4address.address.domain_name
}

resource "qip_v4address_rr" "manually" {
  name        = "extra-example"
  address     = "192.0.2.42"
  domain_name = "int.example.com"
}
