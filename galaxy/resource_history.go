package galaxy

import (
	"context"
	"github.com/brinkmanlab/blend4go"
	"github.com/brinkmanlab/blend4go/histories"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var historyOmitFields = map[string]interface{}{"model_class": nil, "state_ids": nil}

func resourceHistory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHistoryCreate,
		ReadContext:   resourceHistoryRead,
		UpdateContext: resourceHistoryUpdate,
		DeleteContext: resourceHistoryDelete,
		Schema: map[string]*schema.Schema{
			//"id": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			"importable": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Importable",
			},
			"create_time": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time history created",
			},
			"contents_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API url to history contents",
			},
			"size": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total storage size of all containing datasets",
			},
			"user_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User id of assigned user",
			},
			"username_and_slug": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Username and slug",
			},
			"annotation": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Annotation description of history",
			},
			"state_details": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Map of count of datasets keyed on each state",
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Overall state of history and its contents",
			},
			"empty": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "History empty",
			},
			"update_time": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time history last modified",
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of tags assigned to history",
			},
			"deleted": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Deleted",
			},
			"genome_build": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Genome build assigned to history",
			},
			"slug": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Slug",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "History name as displayed to user",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API url of history",
			},
			//"state_ids": &schema.Schema{ TODO flatten?
			//	Type:     schema.TypeMap,
			//	Computed: true,
			//	Elem: &schema.Schema{
			//		Type: schema.TypeList,
			//		Elem: &schema.Schema{
			//			Type: schema.TypeString,
			//		},
			//	},
			//},
			"published": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Published",
			},
			//"model_class": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			"purged": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Purged",
			},
			"purge": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Purge history on delete",
			},
		},
		Description: "Galaxy histories organise and group data into 'workspaces'. All datasets must be associated with a history, including job outputs.",
	}
}

func resourceHistoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	if history, err := histories.NewHistory(ctx, g, d.Get("name").(string)); err == nil {
		return toSchema(history, d, historyOmitFields)
	} else {
		return diag.FromErr(err)
	}
}

func resourceHistoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	if history, err := histories.Get(ctx, g, d.Id()); err == nil {
		return toSchema(history, d, historyOmitFields)
	} else {
		return diag.FromErr(err)
	}
}

func resourceHistoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	history := new(histories.History)
	history.SetGalaxyInstance(g)
	var diags diag.Diagnostics
	diags = append(diags, fromSchema(history, d)...)
	if err := history.Update(ctx); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	diags = append(diags, toSchema(history, d, historyOmitFields)...)
	return diags
}

func resourceHistoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	g := m.(*blend4go.GalaxyInstance)
	history := new(histories.History)
	history.SetGalaxyInstance(g)
	diags = append(diags, fromSchema(history, d)...)
	if err := history.Delete(ctx, d.Get("purge").(bool)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
