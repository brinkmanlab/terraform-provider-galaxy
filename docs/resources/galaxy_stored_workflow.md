# galaxy_stored_workflow Resource

[Galaxy workflows](https://galaxyproject.org/learn/advanced-workflow/) are groups of jobs chained together to process data. This resource represents and manages a workflow stored in a Galaxy instance.

## Example Usage

```hcl
resource "galaxy_stored_workflow" "example" {
  json = file("workflow.ga")
}
```

## Argument Reference

* `annotation` - &lt;String&gt; (Optional) Workflow annotation  
* `import_tools` - &lt;Bool&gt; (Optional) Install tools referenced by workflow  
* `importable` - &lt;Bool&gt; (Optional) Allow users to import workflow  
* `json` - &lt;String&gt; (Required) JSON encoded workflow. See terraform file() to load a .ga file.  
* `name` - &lt;String&gt; (Optional) Name of stored workflow as displayed to user  
* `publish` - &lt;Bool&gt; (Optional) Make workflow available to all users  
* `published` - &lt;Bool&gt; (Optional) Published  
* `show_in_tool_panel` - &lt;Bool&gt; (Optional) Show in tool panel in Galaxy UI  
* `tags` - &lt;List&gt; (Optional) List of tags assigned to workflow  
  Element type: String


## Attribute Reference

* `annotation` - &lt;String&gt; Workflow annotation  
* `deleted` - &lt;Bool&gt; Workflow deleted  
* `import_tools` - &lt;Bool&gt; Install tools referenced by workflow  
* `importable` - &lt;Bool&gt; Allow users to import workflow  
* `json` - &lt;String&gt; JSON encoded workflow. See terraform file() to load a .ga file.  
* `latest_workflow_uuid` - &lt;String&gt; UUID to uniquely identify stored workflow  
* `name` - &lt;String&gt; Name of stored workflow as displayed to user  
* `number_of_steps` - &lt;Int&gt; Count of steps in workflow  
* `owner` - &lt;String&gt; User workflow is assigned to  
* `publish` - &lt;Bool&gt; Make workflow available to all users  
* `published` - &lt;Bool&gt; Published  
* `show_in_tool_panel` - &lt;Bool&gt; Show in tool panel in Galaxy UI  
* `tags` - &lt;List&gt; List of tags assigned to workflow  
  Element type: String
* `url` - &lt;String&gt; URL of workflow within Galaxy API  
* `version` - &lt;Int&gt; Workflow version  

