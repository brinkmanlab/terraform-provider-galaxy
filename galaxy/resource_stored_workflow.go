package galaxy

import (
	"context"
	"github.com/brinkmanlab/blend4go"
	"github.com/brinkmanlab/blend4go/workflows"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var workflowOmitFields = map[string]interface{}{"inputs": nil, "steps": nil, "model_class": nil}

func resourceStoredWorkflow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStoredWorkflowCreate,
		ReadContext:   resourceStoredWorkflowRead,
		UpdateContext: resourceStoredWorkflowUpdate,
		DeleteContext: resourceStoredWorkflowDelete,
		Schema: map[string]*schema.Schema{
			//"id": {
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			"json": {
				Type:        schema.TypeString,
				Required:    true,
				StateFunc:   func(v interface{}) string { return HashString(v.(string)) },
				ForceNew:    true, // https://github.com/galaxyproject/galaxy/issues/10687
				Description: "JSON encoded workflow. See terraform file() to load a .ga file.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of stored workflow as displayed to user",
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true, // TODO https://github.com/galaxyproject/galaxy/issues/10684
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of tags assigned to workflow",
			},
			"deleted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Workflow deleted",
			},
			"latest_workflow_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID to uniquely identify stored workflow",
			},
			"show_in_tool_panel": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Show in tool panel in Galaxy UI",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL of workflow within Galaxy API",
			},
			"number_of_steps": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Count of steps in workflow",
			},
			"published": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true, // TODO https://github.com/galaxyproject/galaxy/issues/10682
				Description: "Make workflow available to all users",
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User workflow is assigned to",
			},
			//"model_class": {
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			//"inputs": {
			//	Type:     schema.TypeList,
			//	Computed: true,
			//},
			"annotation": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Workflow annotation",
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Workflow version",
			},
			//"steps": {
			//	Type:     schema.TypeList,
			//	Computed: true,
			//},
			"import_tools": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Install tools referenced by workflow",
			},
			"importable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true, // TODO https://github.com/galaxyproject/galaxy/issues/10683
				Description: "Allow users to import workflow",
			},
		},
		Importer:    &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Description: "[Galaxy workflows](https://galaxyproject.org/learn/advanced-workflow/) are groups of jobs chained together to process data. This resource represents and manages a workflow stored in a Galaxy instance.",
	}
}

func resourceStoredWorkflowCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	json := d.Get("json").(string)

	if workflow, err := workflows.NewStoredWorkflow(ctx, g, json, d.Get("import_tools").(bool), d.Get("published").(bool), d.Get("importable").(bool)); err == nil {
		return toSchema(workflow, d, workflowOmitFields)
	} else {
		return diag.FromErr(err)
	}
}

func resourceStoredWorkflowRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	if workflow, err := workflows.Get(ctx, g, d.Id()); err == nil {
		return toSchema(workflow, d, workflowOmitFields)
	} else {
		return diag.FromErr(err)
	}
}

func resourceStoredWorkflowUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	var j string
	if d.HasChange("json") {
		j = d.Get("json").(string)
	}

	workflow := new(workflows.StoredWorkflow)
	workflow.SetGalaxyInstance(g)
	workflow.SetID(d.Id())
	workflow.Name = d.Get("name").(string)
	workflow.Annotation = d.Get("annotation").(string)
	workflow.ShowInToolPanel = d.Get("show_in_tool_panel").(bool)

	if err := workflow.Update(ctx, j); err != nil {
		return diag.FromErr(err)
	} else {
		return nil
	}
}

func resourceStoredWorkflowDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	g := m.(*blend4go.GalaxyInstance)
	workflow := new(workflows.StoredWorkflow)
	workflow.SetGalaxyInstance(g)
	diags = append(diags, fromSchema(workflow, d, nil)...)
	if err := workflow.Delete(ctx); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
