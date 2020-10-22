package galaxy

import (
	"context"
	"github.com/brinkmanlab/blend4go"
	"github.com/brinkmanlab/blend4go/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var toolOmitFields = map[string]interface{}{}

func dataSourceTool() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceToolRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the tool to execute in the form `toolshed hostname/repo owner/repo name/tool name/version`",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tool name as displayed to user",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tool version",
			},
			"min_width": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Minimum width",
			},
			"target": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Target",
			},
			"link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Link",
			},
			"panel_section_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tool panel section id",
			},
			"edam_topics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of EDAM Ontology topics associated with tool",
			},
			"form_style": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Form style",
			},
			"edam_operations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of EDAM Ontology operations associated with tool",
			},
			"labels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of labels associated with tool",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tool description",
			},
			"config_file": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Config file",
			},
			"xrefs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Cross references",
			},
			"panel_section_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tool panel section name as displayed to user",
			},
		},
		Description: "Loads information related to an installed tool",
	}
}

func dataSourceToolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	if tool, err := tools.Get(ctx, g, d.Get("id").(string)); err == nil {
		return toSchema(tool, d, toolOmitFields)
	} else {
		return diag.FromErr(err)
	}
}
