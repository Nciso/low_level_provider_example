terraform {
  required_providers {
    provider = {
      version = "0.2.0"
      source  = "company.io/namespace/provider"
    }
  }
}

provider "provider" {
  endpoint = "endpoint"
  token    = "token"
}

data "provider_dummy" "test" {
  dynamic_attribute = "hello"
  regular_attribute = "bye"
  regular_block {
    bar = "bar"
    foo = 3
  }
  dynamic_block {
    bar = "bar"
    foo = 4
  }
}
