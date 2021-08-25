package galaxy

import (
	"context"
	"github.com/brinkmanlab/blend4go"
	"github.com/brinkmanlab/blend4go/repositories"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceToolRepositories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceToolRepositoriesRead,
		Schema: map[string]*schema.Schema{
			"repos": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of repositories",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Repository name",
						},
						"tool_shed": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Toolshed hostname",
						},
						"owner": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Repository owner",
						},
						"changeset_revision": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Changeset revision",
						},
					},
				},
			},
		},
		Description: "Loads list of tool repositories already installed",
	}
}

func dataSourceToolRepositoriesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)
	var diags diag.Diagnostics

	if repos, err := repositories.List(ctx, g); err == nil {
		for _, repo := range repos {
			diags = append(diags, toSchema(repo, d, toolOmitFields)...)
		}
	} else {
		return diag.FromErr(err)
	}
}
