variable "password" {
  type = string
}

resource "galaxy_user" "example" {
  username = "example"
  password = var.password
  email = "example@example.com"
}