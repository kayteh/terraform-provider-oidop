terraform {
  required_providers {
    podio = {
      source = "kayteh/podio"
    }
  }
}

provider "podio" {
  client_id     = var.client_id
  client_secret = var.client_secret
  username      = var.username
  password      = var.password
}
