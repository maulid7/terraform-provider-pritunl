terraform {
  required_providers {
    pritunl = {
      version = "~> 0.0.1"
      source  = "disc/pritunl"
    }
  }
}

provider "pritunl" {
  url    = var.pritunl_url
  token  = var.pritunl_api_token
  secret = var.pritunl_api_secret

  insecure         = var.pritunl_insecure
  connection_check = true
}

resource "pritunl_organization" "developers" {
  name = "Developers"
}

resource "pritunl_organization" "admins" {
  name = "Admins"
}

resource "pritunl_user" "test" {
  name            = "test-user"
  organization_id = pritunl_organization.developers.id
  email           = "test@test.com"
  groups = [
    "admins",
  ]
}

resource "pritunl_user" "test_pin" {
  name            = "test-user-pin"
  organization_id = pritunl_organization.developers.id
  email           = "test@test.com"
  pin             = "123456"
  groups = [
    "admins",
  ]
}

resource "pritunl_server" "test" {
  name = "test"

  organization_ids = [
    pritunl_organization.developers.id,
    pritunl_organization.admins.id,
  ]

  status = "online"
}

resource "pritunl_route" "test" {
  server_id = pritunl_server.test.id

  network     = "8.8.8.8/32"
  comment   = "Google DNS"
  nat       = false 
}

resource "pritunl_route" "test2" {
  server_id = pritunl_server.test.id

  network     = "1.1.1.1/32"
  comment   = "CF DNS"
  nat       = true 
}

resource "pritunl_route" "test3" {
  server_id = pritunl_server.test.id

  network   = "1.2.3.5/32"
  nat       = false 
}
