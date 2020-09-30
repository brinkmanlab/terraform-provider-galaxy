package galaxy

import (
	"context"
	"github.com/brinkmanlab/blend4go"
	"github.com/brinkmanlab/blend4go/workflows"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceStoredWorkflow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStoredWorkflowCreate,
		ReadContext:   resourceStoredWorkflowRead,
		UpdateContext: resourceStoredWorkflowUpdate,
		DeleteContext: resourceStoredWorkflowDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"json": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "JSON encoded workflow. See terraform file() to load a .ga file.",
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"deleted": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"latest_workflow_uuid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"show_in_tool_panel": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"number_of_steps": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"published": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"model_class": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"inputs": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
			},
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"steps": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
			},
		},
	}
}

func resourceStoredWorkflowCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	if workflow, err := workflows.NewStoredWorkflow(ctx, g, d.Get("json").(string)); err == nil {
		return toSchema(workflow, d)
	} else {
		return diag.FromErr(err)
	}
}

func resourceStoredWorkflowRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	if workflow, err := workflows.Get(ctx, g, d.Get("id").(string)); err == nil {
		if json, err := workflow.Download(ctx); err == nil {
			d.Set("json", json)
		} else {
			return diag.FromErr(err)
		}
		return toSchema(workflow, d)
	} else {
		return diag.FromErr(err)
	}
}

func resourceStoredWorkflowUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	workflow := new(workflows.StoredWorkflow)
	workflow.SetGalaxyInstance(g)
	workflow.SetID(d.Get("id").(string))
	workflow.Name = d.Get("name").(string)
	workflow.Annotation = d.Get("annotation").(string)
	workflow.ShowInToolPanel = d.Get("show_in_tool_panel").(bool)
	return diag.FromErr(workflow.Update(ctx, d.Get("json").(string)))
}

func resourceStoredWorkflowDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	g := m.(*blend4go.GalaxyInstance)
	workflow := new(workflows.StoredWorkflow)
	workflow.SetGalaxyInstance(g)
	diags = append(diags, fromSchema(workflow, d)...)
	diags = append(diags, diag.FromErr(workflow.Delete(ctx))...)

	return diags
}
