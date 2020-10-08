package galaxy

import (
	"context"
	"github.com/brinkmanlab/blend4go"
	"github.com/brinkmanlab/blend4go/histories"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceHistory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHistoryCreate,
		ReadContext:   resourceHistoryRead,
		UpdateContext: resourceHistoryUpdate,
		DeleteContext: resourceHistoryDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"importable": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"create_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"contents_url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"user_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"username_and_slug": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"state_details": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"empty": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"update_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"deleted": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"genome_build": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"slug": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"state_ids": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			"published": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"model_class": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"purged": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"purge": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func historyUpdate(ctx context.Context, history *histories.History, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	diags = append(diags, fromSchema(history, d)...)
	diags = append(diags, diag.FromErr(history.Update(ctx))...)
	diags = append(diags, toSchema(history, d)...)
	return diags
}

func resourceHistoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	if history, err := histories.NewHistory(ctx, g, d.Get("name").(string)); err == nil {
		return historyUpdate(ctx, history, d)
	} else {
		return diag.FromErr(err)
	}
}

func resourceHistoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	if history, err := histories.Get(ctx, g, d.Get("id").(string)); err == nil {
		return toSchema(history, d)
	} else {
		return diag.FromErr(err)
	}
}

func resourceHistoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	history := new(histories.History)
	history.SetGalaxyInstance(g)
	return historyUpdate(ctx, history, d)
}

func resourceHistoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	g := m.(*blend4go.GalaxyInstance)
	history := new(histories.History)
	history.SetGalaxyInstance(g)
	diags = append(diags, fromSchema(history, d)...)
	diags = append(diags, diag.FromErr(history.Delete(ctx, d.Get("purge").(bool)))...)

	return diags
}
