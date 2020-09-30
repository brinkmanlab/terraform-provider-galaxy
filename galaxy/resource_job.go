package galaxy

import (
	"context"
	"fmt"
	"github.com/brinkmanlab/blend4go"
	"github.com/brinkmanlab/blend4go/jobs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
)

func resourceJob() *schema.Resource {
	job := map[string]*schema.Schema{
		"id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"tool_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"update_time": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"history_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"exit_code": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"state": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"create_time": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"model_class": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		//"inputs": &schema.Schema{
		//	Type:     interface{},
		//	Computed: true,
		//},
		//"outputs": &schema.Schema{
		//	Type:     interface{},
		//	Computed: true,
		//},
		//"params": &schema.Schema{
		//	Type:     interface{},
		//	Computed: true,
		//},
	}

	firstJob := map[string]*schema.Schema{
		"tool_id": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			ExactlyOneOf: []string{"tool_id", "tool_uuid"},
		},
		"tool_uuid": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			ExactlyOneOf: []string{"tool_id", "tool_uuid"},
		},
		"history_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"param": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"hda": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"hdca": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"additional_jobs": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				ReadContext:   resourceJobRead,
				DeleteContext: resourceJobDelete,
				Schema:        job,
			},
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
	}
}

func jobsToSchema(job []*jobs.Job, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	diags = append(diags, toSchema(job[0], d)...)
	var additionalJobs []map[string]interface{}
	for _, j := range job[1:] {
		aj := map[string]interface{}{}
		if decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{TagName: "json", Result: aj}); err == nil {
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
	if id, ok := d.GetOk("tool_uuid"); ok {
		payload["tool_uuid"] = id.(string)
	}

	// Prepare param inputs
	if params, ok := d.GetOk("param"); ok {
		for i, param := range params.(map[string]map[string]string) {
			name, name_ok := param["name"]
			if !name_ok || name == "" {
				diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: fmt.Sprintf("Required 'name' field missing in param #%s", i)})
			}
			value, value_ok := param["value"]
			if !value_ok {
				diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: fmt.Sprintf("Required 'value' field missing in param #%s", i)})
			}
			if name_ok && name != "" && value_ok {
				inputs[name] = value
			}
		}
	}

	// Prepare hda and hdca inputs
	for _, t := range []string{"hda", "hdca"} {
		if hdas, ok := d.GetOk(t); ok {
			for i, hda := range hdas.(map[string]map[string]string) {
				name, name_ok := hda["name"]
				if !name_ok || name == "" {
					diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: fmt.Sprintf("Required 'name' field missing in %s #%s", t, i)})
				}
				value, value_ok := hda["value"]
				if !value_ok {
					diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: fmt.Sprintf("Required 'value' field missing in %s #%s", t, i)})
				}
				if name_ok && name != "" && value_ok {
					in, ok := inputs[name] // All hda with the same input name are merged into same list
					input := in.(*map[string][]map[string]string)
					if !ok {
						input = &map[string][]map[string]string{}
						inputs[name] = input
					}
					(*input)["values"] = append((*input)["values"], map[string]string{"id": value, "src": t})
					// TODO 'batch'?
				}
			}
		}
	}

	if job, _, _, _, err := jobs.NewJob(ctx, g, payload); err == nil { //TODO Expose job outputs?
		diags = append(diags, jobsToSchema(job, d)...)
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
	if job, err := jobs.Get(ctx, g, d.Get("id").(string)); err == nil {
		jobList = append(jobList, job)
	} else {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Get additional jobs
	for _, job := range d.Get("additional_jobs").([]map[string]interface{}) {
		if job, err := jobs.Get(ctx, g, job["id"].(string)); err == nil {
			jobList = append(jobList, job)
		} else {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Don't apply to schema if errors
	if diags.HasError() {
		return diags
	}
	return jobsToSchema(jobList, d)
}

func resourceJobDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	g := m.(*blend4go.GalaxyInstance)
	job := new(jobs.Job)
	job.SetGalaxyInstance(g)
	diags = append(diags, fromSchema(job, d)...)
	diags = append(diags, diag.FromErr(job.Delete(ctx))...)

	return diags
}
