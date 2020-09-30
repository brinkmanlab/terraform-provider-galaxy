package galaxy

import (
	"context"
	"github.com/brinkmanlab/blend4go"
	"github.com/brinkmanlab/blend4go/repositories"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRepository() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRepositoryCreate,
		ReadContext:   resourceRepositoryRead,
		DeleteContext: resourceRepositoryDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"deleted": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ctx_rev": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_message": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"installed_changeset_revision": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tool_shed": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"dist_to_shed": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"uninstalled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"changeset_revision": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"include_datatypes": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"latest_installable_revision": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"revision_update": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"revision_upgrade": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"repository_deprecated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"install_tool_dependencies": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"install_repository_dependencies": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"install_resolver_dependencies": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tool_panel_section_id": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"tool_panel_section_id", "new_tool_panel_section_label"},
			},
			"new_tool_panel_section_label": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "",
				ConflictsWith: []string{"tool_panel_section_id", "new_tool_panel_section_label"},
			},
			"remove_from_disk": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceRepositoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	if err := repositories.Install(ctx, g,
		d.Get("tool_shed").(string),
		d.Get("name").(string),
		d.Get("owner").(string),
		d.Get("changeset_revision").(string),
		d.Get("install_tool_dependencies").(bool),
		d.Get("install_repository_dependencies").(bool),
		d.Get("install_resolver_dependencies").(bool),
		d.Get("tool_panel_section_id").(string),
		d.Get("new_tool_panel_section_label").(string),
	); err == nil {
		return nil //toSchema(workflow, d)
	} else {
		return diag.FromErr(err)
	}
}

func resourceRepositoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	if repo, err := repositories.Get(ctx, g, d.Get("id").(string)); err == nil {
		return toSchema(repo, d)
	} else {
		return diag.FromErr(err)
	}
}

func resourceRepositoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	if err := repositories.UninstallID(ctx, g, d.Get("id").(string), d.Get("remove_from_disk").(bool)); err == nil {
		return nil
	} else {
		return diag.FromErr(err)
	}
}
