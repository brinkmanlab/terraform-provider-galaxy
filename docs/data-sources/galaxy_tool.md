# galaxy_tool Data Source



## Example Usage

```hcl
resource "galaxy_repository" "awkscript" {
  tool_shed = "toolshed.g2.bx.psu.edu"
  owner = "brinkmanlab"
  name = "awkscript"
  changeset_revision = "ceac6ffb3865"
  remove_from_disk = true
}

data "galaxy_tool" "awkscript" {
  depends_on = [galaxy_repository.awkscript]
  id = "toolshed.g2.bx.psu.edu/repos/brinkmanlab/awkscript/awkscript/1.0"  # An issue with the Galaxy API requires this. https://github.com/galaxyproject/galaxy/issues/10378
}
```

## Argument Reference

* `id` - &lt;String&gt; (Required) Id of the tool to execute in the form `toolshed hostname/repo owner/repo name/tool name/version`  


## Attribute Reference

* `config_file` - &lt;String&gt; Config file  
* `description` - &lt;String&gt; Tool description  
* `edam_operations` - &lt;List&gt; List of EDAM Ontology operations associated with tool  
  Element type: String
* `edam_topics` - &lt;List&gt; List of EDAM Ontology topics associated with tool  
  Element type: String
* `form_style` - &lt;String&gt; Form style  
* `id` - &lt;String&gt; Id of the tool to execute in the form `toolshed hostname/repo owner/repo name/tool name/version`  
* `labels` - &lt;List&gt; List of labels associated with tool  
  Element type: String
* `link` - &lt;String&gt; Link  
* `min_width` - &lt;Int&gt; Minimum width  
* `name` - &lt;String&gt; Tool name as displayed to user  
* `panel_section_id` - &lt;String&gt; Tool panel section id  
* `panel_section_name` - &lt;String&gt; Tool panel section name as displayed to user  
* `target` - &lt;String&gt; Target  
* `version` - &lt;String&gt; Tool version  
* `xrefs` - &lt;List&gt; Cross references  
  Element type: String

