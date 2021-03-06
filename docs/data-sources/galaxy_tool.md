# galaxy_tool Data Source

Loads information related to an installed tool

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
  id = galaxy_repository.awkscript.tools[0].tool_guid
}
```

## Argument Reference

* `guid` - &lt;String&gt; (Optional) Tool guid  
  Exactly one of `id` or `guid`  
* `id` - &lt;String&gt; (Optional) Tool Id  
  Exactly one of `id` or `guid`  


## Attribute Reference

* `config_file` - &lt;String&gt; Config file  
* `description` - &lt;String&gt; Tool description  
* `edam_operations` - &lt;List&gt; List of EDAM Ontology operations associated with tool  
  Element type: String
* `edam_topics` - &lt;List&gt; List of EDAM Ontology topics associated with tool  
  Element type: String
* `form_style` - &lt;String&gt; Form style  
* `guid` - &lt;String&gt; Tool guid  
* `id` - &lt;String&gt; Tool Id  
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

