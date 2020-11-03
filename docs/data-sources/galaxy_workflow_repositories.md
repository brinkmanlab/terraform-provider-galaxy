# galaxy_workflow_repositories Data Source

Galaxy workflows are dependant on the presence of the tools they use to be installed in the same Galaxy instance. The data source extracts the tool repositories referenced within the workflow json, listing them for installation. See [resource_repository](../resources/resource_repository) for more information on installing repositories.

## Example Usage

```hcl
data "galaxy_workflow_repositories" "example" {
  json = file("workflow.ga")
}

resource "galaxy_repository" "example" {
  for_each = data.galaxy_workflow_repositories.example.repositories
  tool_shed = each.value.tool_shed
  owner = each.value.owner
  name = each.value.name
  changeset_revision = each.value.changeset_revision
  remove_from_disk = true
}
```

## Argument Reference

* `json` - &lt;String&gt; (Required) JSON encoded workflow. See terraform file() to load a .ga file.  


## Attribute Reference

* `json` - &lt;String&gt; JSON encoded workflow. See terraform file() to load a .ga file.  
* `repositories` - &lt;Set&gt; Set of repositories referenced within workflow  
  Attributes:  
  * `changeset_revision` - &lt;String&gt; Changeset revision  
  * `name` - &lt;String&gt; Repository name  
  * `owner` - &lt;String&gt; Repository owner  
  * `tool_shed` - &lt;String&gt; Toolshed hostname  


