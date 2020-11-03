---
page_title: Bootstrap a Galaxy admin API key
---
# Bootstrap a Galaxy admin API key using the master API key

Many operations are not permitted using Galaxies master API key. 
To properly configure this provider you need to use an administrator API key.
This assumes you have configured `admin_users: admin@example.com` in the Galaxy configuration.


```hcl
locals {
  host = "http://localhost:8080"
}

variable "password" {
  type = string
}

# Configure a aliased Galaxy provider with the master API key
provider "galaxy" {
  host = local.host
  api_key = "master API key"
  alias = "master"
}

# Use the master API key to create the admin user
resource "galaxy_user" "admin" {
  provider = galaxy.master
  username = "admin"
  password = var.password
  email = "admin@example.com"
}

# Configure the primary Galaxy provider to use the admin API key
provider "galaxy" {
  host = local.host
  api_key = galaxy_user.admin.api_key
}
```
