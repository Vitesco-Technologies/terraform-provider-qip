resource "qip_v4address" "simple" {
  address = "192.0.2.23"
  subnet  = "192.0.2.0"
  name    = "my-example"
}

resource "qip_v4address" "full" {
  address = "192.0.2.23"
  subnet  = "192.0.2.0"

  subnet_range_start = "192.0.2.30"
  subnet_range_end   = "192.0.2.50"

  name         = "my-example"
  object_class = "Virtual Server"
  description  = "Example System"
  domain_name  = "corp.example.com"
}
