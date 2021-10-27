package galaxy

import (
	"context"
	"errors"
	"github.com/brinkmanlab/blend4go"
	"github.com/brinkmanlab/blend4go/jobs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
	"time"
)

var jobOmitFields = map[string]interface{}{"inputs": nil, "outputs": nil, "params": nil, "model_class": nil}
var jobEnded = map[string]bool{
	"new":         false,
	"upload":      false,
	"waiting":     false,
	"queued":      false,
	"running":     false,
	"ok":          true,
	"error":       true,
	"paused":      true,
	"deleted":     true,
	"deleted_new": true,
}

func resourceJob() *schema.Resource {
	job := map[string]*schema.Schema{
		//"id": {
		//	Type:     schema.TypeString,
		//	Computed: true,
		//},
		"tool_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Id of the tool to execute in the form `toolshed hostname/repo owner/repo name/tool name/version`",
		},
		"update_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Time job state lst updated",
		},
		"history_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Id of history where tool outputs are associated",
		},
		"exit_code": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Exit code as returned by tool execution",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Running state of job",
		},
		"create_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Job creation time",
		},
		//"model_class": {
		//	Type:     schema.TypeString,
		//	Computed: true,
		//},
		//"inputs": {
		//	Type:     interface{},
		//	Computed: true,
		//},
		//"outputs": {
		//	Type:     interface{},
		//	Computed: true,
		//},
		//"params": {
		//	Type:     interface{},
		//	Computed: true,
		//},
	}

	firstJob := map[string]*schema.Schema{
		"tool_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ExactlyOneOf: []string{"tool_id", "tool_guid"},
			ForceNew:     true,
			Description:  "Id of the tool to execute in the form `toolshed hostname/repos/repo owner/repo name/tool name/version`",
		},
		"tool_guid": {
			Type:         schema.TypeString,
			Optional:     true,
			ExactlyOneOf: []string{"tool_id", "tool_guid"},
			ForceNew:     true,
			Description:  "UUID of tool as assigned by Galaxy instance",
		},
		"history_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Id of history where tool outputs are associated",
		},
		"params": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			ForceNew:    true,
			Description: "Map of parameter values keyed on input id",
		},
		"hda": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"input": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Input id as described in Galaxy tool wrapper XML",
					},
					"id": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "HDA id",
					},
				},
			},
			ForceNew:    true,
			Description: "Repeatable block of HDA inputs. Specify the same input id in multiple blocks to provide tool multiple HDAs per input.",
		},
		"hdca": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"input": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Input id as described in Galaxy tool wrapper XML",
					},
					"id": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "HDCA id",
					},
				},
			},
			ForceNew:    true,
			Description: "Repeatable block of HDCA inputs. Specify the same input id in multiple blocks to provide tool multiple HDCAs per input.",
		},
		"additional_jobs": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				ReadContext:   resourceJobRead,
				DeleteContext: resourceJobDelete,
				Schema:        job,
			},
			Description: "If the input parameters spawn multiple jobs, the remaining jobs will be listed here",
		},
		"wait_for_completion": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			ForceNew:    true,
			Description: "Wait for job to complete before creating dependant resources",
		},
	}

	// Copy in common job fields
	for k, v := range job {
		if _, ok := firstJob[k]; !ok { // Do not overwrite redeclared fields
			firstJob[k] = v
		}
	}

	return &schema.Resource{
		CreateContext: resourceJobCreate,
		ReadContext:   resourceJobRead,
		DeleteContext: resourceJobDelete,
		Schema:        firstJob,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Description:   "Execute tools to load data. This is mainly intended for data managers or upload/download tools. Do not use this for data processing!",
	}
}

// Handle Job array, converting the first job to the main schema and all others assigned to the additional_jobs field
func jobsToSchema(job []*jobs.Job, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	diags = append(diags, toSchema(job[0], d, jobOmitFields)...)
	var additionalJobs []map[string]interface{}
	for _, j := range job[1:] {
		aj := map[string]interface{}{}
		if decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{TagName: "json", Result: &aj}); err == nil {
			if err := decoder.Decode(j); err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		} else {
			diags = append(diags, diag.FromErr(err)...)
		}
		additionalJobs = append(additionalJobs, aj)
	}
	if err := d.Set("additional_jobs", additionalJobs); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	return diags
}

func resourceJobCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	g := m.(*blend4go.GalaxyInstance)
	inputs := map[string]interface{}{}
	payload := map[string]interface{}{
		"history_id": d.Get("history_id"),
		"inputs":     &inputs,
	}
	if id, ok := d.GetOk("tool_id"); ok {
		payload["tool_id"] = id.(string)
	}
	if id, ok := d.GetOk("tool_guid"); ok {
		payload["tool_uuid"] = id.(string)
	}

	// Prepare param inputs
	if params, ok := d.GetOk("params"); ok {
		for k, v := range params.(map[string]interface{}) {
			inputs[k] = v.(string)
		}
	}

	// Prepare hda and hdca inputs
	for _, t := range []string{"hda", "hdca"} {
		if hdas, ok := d.GetOk(t); ok {
			for _, hda := range hdas.([]map[string]string) {
				name := hda["input"]
				in, ok := inputs[name] // All hda with the same input name are merged into same list
				var input *map[string][]map[string]string
				if ok {
					input = in.(*map[string][]map[string]string)
				} else {
					input = &map[string][]map[string]string{}
					inputs[name] = input
				}
				(*input)["values"] = append((*input)["values"], map[string]string{"id": hda["id"], "src": t})
				// TODO 'batch'?
			}
		}
	}

	if jobList, _, _, _, err := jobs.NewJob(ctx, g, payload); err == nil { //TODO Expose job outputs?
		if d.Get("wait_for_completion").(bool) {
			complete := false
			for !complete {
				complete = true
				time.Sleep(2 * time.Second)
				for _, job := range jobList {
					if _, err := g.Get(ctx, job.GetID(), job, &map[string]string{}); err == nil {
						complete = complete && jobEnded[job.State]
					} else {
						diags = append(diags, diag.FromErr(err)...)
						complete = false
					}
				}
				if errors.Is(ctx.Err(), context.Canceled) {
					complete = true
					break
				}
			}

			// If waiting on jobs, that means we need them to succeed
			for _, job := range jobList {
				if job.State == "error" {
					diags = append(diags, diag.Errorf("galaxy_job failed execution, see %v/api/jobs/%v for more info", g.Client.HostURL, job.Id)...)
				}
			}
		}
		diags = append(diags, jobsToSchema(jobList, d)...)
	} else {
		return diag.FromErr(err)
	}
	return diags
}

func resourceJobRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	g := m.(*blend4go.GalaxyInstance)
	var jobList []*jobs.Job

	// Get first job
	if job, err := jobs.Get(ctx, g, d.Id()); err == nil {
		jobList = append(jobList, job)
	} else {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Get additional jobs
	for _, job := range d.Get("additional_jobs").([]interface{}) {
		if job, err := jobs.Get(ctx, g, job.(map[string]interface{})["id"].(string)); err == nil {
			jobList = append(jobList, job)
		} else {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Don't apply to schema if errors
	if diags != nil && diags.HasError() {
		return diags
	}
	return jobsToSchema(jobList, d)
}

func resourceJobDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	g := m.(*blend4go.GalaxyInstance)
	job := new(jobs.Job)
	job.SetGalaxyInstance(g)
	diags = append(diags, fromSchema(job, d, nil)...)
	if err := job.Delete(ctx); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
