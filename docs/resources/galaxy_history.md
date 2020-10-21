# galaxy_history Resource

Galaxy histories organise and group data into &#39;workspaces&#39;. All datasets must be associated with a history, including job outputs.

## Example Usage

```hcl
resource "galaxy_history" "example" {
  name = "example"
}
```

## Argument Reference

* `annotation` - &lt;String&gt; (Optional) Annotation description of history  
* `genome_build` - &lt;String&gt; (Optional) Genome build assigned to history  
* `name` - &lt;String&gt; (Optional) History name as displayed to user  
* `published` - &lt;Bool&gt; (Optional) Published  
* `purge` - &lt;Bool&gt; (Optional) Purge history on delete \[Default: true]  
* `slug` - &lt;String&gt; (Optional) Slug  
* `tags` - &lt;List&gt; (Optional) List of tags assigned to history  
  Element type: String


## Attribute Reference

* `annotation` - &lt;String&gt; Annotation description of history  
* `contents_url` - &lt;String&gt; API url to history contents  
* `create_time` - &lt;String&gt; Time history created  
* `deleted` - &lt;Bool&gt; Deleted  
* `empty` - &lt;Bool&gt; History empty  
* `genome_build` - &lt;String&gt; Genome build assigned to history  
* `importable` - &lt;Bool&gt; Importable  
* `name` - &lt;String&gt; History name as displayed to user  
* `published` - &lt;Bool&gt; Published  
* `purge` - &lt;Bool&gt; Purge history on delete  
* `purged` - &lt;Bool&gt; Purged  
* `size` - &lt;Int&gt; Total storage size of all containing datasets  
* `slug` - &lt;String&gt; Slug  
* `state` - &lt;String&gt; Overall state of history and its contents  
* `state_details` - &lt;Map&gt; Map of count of datasets keyed on each state  
  Element type: Int
* `tags` - &lt;List&gt; List of tags assigned to history  
  Element type: String
* `update_time` - &lt;String&gt; Time history last modified  
* `url` - &lt;String&gt; API url of history  
* `user_id` - &lt;String&gt; User id of assigned user  
* `username_and_slug` - &lt;String&gt; Username and slug  

