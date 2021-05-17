variable "password" {
  type = string
}

resource "galaxy_user" "example" {
  username = "example"
  password = var.password
  email = "example@example.com"
}

resource "galaxy_quota" "example" {
  name = "example"
  description = "example quota"
  amount = "1G"
  users = [galaxy_user.example.id]
}