package galaxy

import (
	"context"
	"github.com/brinkmanlab/blend4go"
	"github.com/brinkmanlab/blend4go/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTool() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceToolRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_width": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"target": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"panel_section_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"edam_topics": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"form_style": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"edam_operations": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"labels": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"config_file": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"xrefs": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"panel_section_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceToolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	if tool, err := tools.Get(ctx, g, d.Get("id").(string)); err == nil {
		return toSchema(tool, d)
	} else {
		return diag.FromErr(err)
	}
}
