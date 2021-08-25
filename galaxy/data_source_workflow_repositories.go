package galaxy

import (
	"context"
	"github.com/brinkmanlab/blend4go/workflows"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceWorkflowRepositories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWorkflowRepositoriesRead,
		Schema: map[string]*schema.Schema{
			"json": {
				Type:        schema.TypeString,
				Required:    true,
				StateFunc:   func(v interface{}) string { return HashString(v.(string)) },
				Description: "JSON encoded workflow. See terraform file() to load a .ga file.",
			},
			"repositories": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Refactor resource_repository fields and import here
					},
				},
				Description: "Set of repositories referenced within workflow",
			},
		},
		Description: "Galaxy workflows are dependant on the presence of the tools they use to be installed in the same Galaxy instance. The data source extracts the tool repositories referenced within the workflow json, listing them for installation. See [resource_repository](../resources/resource_repository) for more information on installing repositories.",
	}
}

func dataSourceWorkflowRepositoriesRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	json := d.Get("json").(string)
	hash := HashString(json)
	if repos, err := workflows.Repositories(json); err == nil {
		r := make([]map[string]string, len(repos))
		for i, repo := range repos {
			r[i] = map[string]string{
				"name":               repo.Name,
				"tool_shed":          repo.ToolShed,
				"owner":              repo.Owner,
				"changeset_revision": repo.ChangesetRevision,
			}
		}
		if err := d.Set("repositories", r); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(hash)
	} else {
		return diag.FromErr(err)
	}
	return nil
}
