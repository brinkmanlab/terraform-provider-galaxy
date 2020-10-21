# galaxy_user Resource

Create and manage Galaxy users. Used mostly to configure admin users.

## Example Usage

```hcl
variable "password" {
  type = string
}

resource "galaxy_user" "example" {
  username = "example"
  password = var.password
  email = "example@example.com"
}
```

## Argument Reference

* `email` - &lt;String&gt; (Required) Users email address  
* `password` - &lt;String&gt; (Required) Password to authenticate user against Galaxy  
* `username` - &lt;String&gt; (Required) Username to identify user  


## Attribute Reference

* `api_key` - &lt;String&gt; API key of user  
* `deleted` - &lt;Bool&gt; User deleted  
* `email` - &lt;String&gt; Users email address  
* `is_admin` - &lt;Bool&gt; User is administrator  
* `nice_total_disk_usage` - &lt;String&gt; Human readable total disk usage of users stored data  
* `password` - &lt;String&gt; Password to authenticate user against Galaxy  
* `purged` - &lt;Bool&gt; User purged  
* `quota` - &lt;String&gt; Maximum disk storage available to user  
* `quota_percent` - &lt;Int&gt; Storage quota, between 0 and 100  
* `tags_used` - &lt;List&gt; List of tags assigned to users resources  
  Element type: String
* `total_disk_usage` - &lt;Float&gt; Total disk usage of users stored data  
* `username` - &lt;String&gt; Username to identify user  

