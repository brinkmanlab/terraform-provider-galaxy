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
			"json": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "JSON encoded workflow. See terraform file() to load a .ga file.",
			},
			"repositories": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resourceRepository(),
			},
		},
	}
}

func dataSourceWorkflowRepositoriesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if repos, err := workflows.Repositories(d.Get("json").(string)); err == nil {
		d.Set("repositories", repos)
	} else {
		return diag.FromErr(err)
	}
	return nil
}
