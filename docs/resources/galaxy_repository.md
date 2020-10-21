# galaxy_repository Resource



## Example Usage

```hcl
resource "galaxy_repository" "example" {
  tool_shed = "toolshed.g2.bx.psu.edu"
  owner = "brinkmanlab"
  name = "awkscript"
  changeset_revision = "ceac6ffb3865"
  remove_from_disk = true
}
```

## Argument Reference

* `changeset_revision` - &lt;String&gt; (Required) Changeset revision of repository  
* `install_repository_dependencies` - &lt;Bool&gt; *Depreciated* (Optional) Install repository dependencies from toolshed  
* `install_resolver_dependencies` - &lt;Bool&gt; (Optional) Install resolver dependencies  
* `install_tool_dependencies` - &lt;Bool&gt; (Optional) Install tool dependencies using the configured dependency manager  
* `name` - &lt;String&gt; (Required) Repository name  
* `new_tool_panel_section_label` - &lt;String&gt; (Optional) Label of tool panel section to create and list tool under  
  Conflicts with `tool_panel_section_id`  
* `owner` - &lt;String&gt; (Required) Repository owner  
* `remove_from_disk` - &lt;Bool&gt; (Optional) Repository files from disk on uninstall \[Default: true]  
* `tool_panel_section_id` - &lt;String&gt; (Optional) Tool panel section ID to list tool under  
  Conflicts with `new_tool_panel_section_label`  
* `tool_shed` - &lt;String&gt; (Required) Repository toolshed  


## Attribute Reference

* `changeset_revision` - &lt;String&gt; Changeset revision of repository  
* `ctx_rev` - &lt;String&gt;   
* `deleted` - &lt;Bool&gt; Repository deleted  
* `dist_to_shed` - &lt;Bool&gt;   
* `error_message` - &lt;String&gt; Install error message  
* `include_datatypes` - &lt;Bool&gt; Repository includes datatypes  
* `install_repository_dependencies` - &lt;Bool&gt; *Depreciated* Install repository dependencies from toolshed  
* `install_resolver_dependencies` - &lt;Bool&gt; Install resolver dependencies  
* `install_tool_dependencies` - &lt;Bool&gt; Install tool dependencies using the configured dependency manager  
* `installed_changeset_revision` - &lt;String&gt; Installed changeset revision  
* `latest_installable_revision` - &lt;String&gt; Latest installable revision of repository  
* `name` - &lt;String&gt; Repository name  
* `new_tool_panel_section_label` - &lt;String&gt; Label of tool panel section to create and list tool under  
* `owner` - &lt;String&gt; Repository owner  
* `remove_from_disk` - &lt;Bool&gt; Repository files from disk on uninstall  
* `repository_deprecated` - &lt;String&gt; Repository depreciated. https://github.com/galaxyproject/galaxy/issues/10453  
* `revision_update` - &lt;String&gt; https://github.com/galaxyproject/galaxy/issues/10453  
* `revision_upgrade` - &lt;String&gt; https://github.com/galaxyproject/galaxy/issues/10453  
* `status` - &lt;String&gt; Installation status  
* `tool_panel_section_id` - &lt;String&gt; Tool panel section ID to list tool under  
* `tool_shed` - &lt;String&gt; Repository toolshed  
* `uninstalled` - &lt;Bool&gt; Uninstalled  
* `url` - &lt;String&gt; Repository url  
