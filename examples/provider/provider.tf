terraform {
  required_providers {
    webdock = {
      source = "hashicorp.com/edu/webdock"
    }
  }
  required_version = ">= 1.1.0"
}


provider "webdock" {
  token = ""
}
