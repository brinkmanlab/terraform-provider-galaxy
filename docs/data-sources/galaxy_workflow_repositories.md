# galaxy_workflow_repositories Data Source



## Example Usage

```hcl
data "galaxy_workflow_repositories" "example" {
  json = file("workflow.ga")
}
```

## Argument Reference

* `json` - &lt;String&gt; (Required) JSON encoded workflow. See terraform file() to load a .ga file.  


## Attribute Reference

* `json` - &lt;String&gt; JSON encoded workflow. See terraform file() to load a .ga file.  
* `repositories` - &lt;List&gt;   
  Attributes:  
  * `changeset_revision` - &lt;String&gt; Changeset revision  
  * `name` - &lt;String&gt; Repository name  
  * `owner` - &lt;String&gt; Repository owner  
  * `tool_shed` - &lt;String&gt; Toolshed hostname  


