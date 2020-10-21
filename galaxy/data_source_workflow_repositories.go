package galaxy

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"github.com/brinkmanlab/blend4go/workflows"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceWorkflowRepositories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWorkflowRepositoriesRead,
		Schema: map[string]*schema.Schema{
			"json": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "JSON encoded workflow. See terraform file() to load a .ga file.",
			},
			"repositories": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Repository name",
						},
						"tool_shed": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Toolshed hostname",
						},
						"owner": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Repository owner",
						},
						"changeset_revision": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Changeset revision",
						},
					},
				},
				Description: "Set of repositories referenced within workflow",
			},
		},
		Description: "Galaxy workflows are dependant on the presence of the tools they use to be installed in the same Galaxy instance. The data source extracts the tool repositories referenced within the workflow json, listing them for installation. See [resource_repository](./resource_repository) for more information on installing repositories.",
	}
}

func dataSourceWorkflowRepositoriesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if repos, err := workflows.Repositories(d.Get("json").(string)); err == nil {
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

		// Use sha1 of json as ID
		h := sha1.New()
		h.Write([]byte(d.Get("json").(string)))
		d.SetId(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	} else {
		return diag.FromErr(err)
	}
	return nil
}
