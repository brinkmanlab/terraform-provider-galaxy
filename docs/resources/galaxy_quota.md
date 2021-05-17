# galaxy_quota Resource

Galaxy quotas regulate the amount of data a user can store in their account at any one time.

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

resource "galaxy_quota" "example" {
  name = "example"
  description = "example quota"
  amount = "1G"
  users = [galaxy_user.example.id]
}
```

## Argument Reference

* `amount` - &lt;String&gt; (Required) Examples: &#34;10000MB&#34;, &#34;99 gb&#34;, &#34;0.2T&#34;, &#34;unlimited&#34;  
* `default` - &lt;String&gt; (Optional) Set as default for category of users (unregistered, registered) \[Default: no]  
* `description` - &lt;String&gt; (Required) Description of quota  
* `groups` - &lt;List&gt; (Optional) List of group ids to apply quota to  
  At least one of `users`, `groups` or `default`  
  Element type: String
* `name` - &lt;String&gt; (Optional) Quota name as displayed to user  
* `operation` - &lt;String&gt; (Optional) Assign (=), increase by amount (+), or decrease by amount (-) \[Default: =]  
* `purge` - &lt;Bool&gt; (Optional) Purge a user on deletion \[Default: true]  
* `users` - &lt;List&gt; (Optional) List of user ids to apply quota to  
  At least one of `users`, `groups` or `default`  
  Element type: String


## Attribute Reference

* `amount` - &lt;String&gt; Examples: &#34;10000MB&#34;, &#34;99 gb&#34;, &#34;0.2T&#34;, &#34;unlimited&#34;  
* `bytes` - &lt;Int&gt; Amount, in bytes  
* `default` - &lt;String&gt; Set as default for category of users (unregistered, registered)  
* `description` - &lt;String&gt; Description of quota  
* `display_amount` - &lt;String&gt; Human readable amount  
* `groups` - &lt;List&gt; List of group ids to apply quota to  
  Element type: String
* `name` - &lt;String&gt; Quota name as displayed to user  
* `operation` - &lt;String&gt; Assign (=), increase by amount (+), or decrease by amount (-)  
* `purge` - &lt;Bool&gt; Purge a user on deletion  
* `users` - &lt;List&gt; List of user ids to apply quota to  
  Element type: String

