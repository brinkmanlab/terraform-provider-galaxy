package galaxy

import (
	"context"
	"fmt"
	"github.com/brinkmanlab/blend4go"
	"github.com/brinkmanlab/blend4go/repositories"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var repositoryOmitFields = map[string]interface{}{"tool_shed_status": nil}

func resourceRepository() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRepositoryCreate,
		ReadContext:   resourceRepositoryRead,
		UpdateContext: resourceRepositoryUpdate,
		DeleteContext: resourceRepositoryDelete,
		Schema: map[string]*schema.Schema{
			//"id": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Installation status",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Repository name",
			},
			"deleted": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Repository deleted",
			},
			"ctx_rev": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_message": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Install error message",
			},
			"installed_changeset_revision": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Installed changeset revision",
			},
			"tool_shed": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Repository toolshed",
			},
			"dist_to_shed": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Repository url",
			},
			"uninstalled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Uninstalled",
			},
			"owner": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Repository owner",
			},
			"changeset_revision": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Changeset revision of repository",
			},
			"include_datatypes": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Repository includes datatypes",
			},
			"latest_installable_revision": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Latest installable revision of repository",
			},
			"revision_update": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "https://github.com/galaxyproject/galaxy/issues/10453",
			},
			"revision_upgrade": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "https://github.com/galaxyproject/galaxy/issues/10453",
			},
			"repository_deprecated": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Repository depreciated. https://github.com/galaxyproject/galaxy/issues/10453",
			},
			"install_tool_dependencies": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Install tool dependencies using the configured dependency manager",
			},
			"install_repository_dependencies": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Deprecated:  "Repository dependencies are depreciated",
				Description: "Install repository dependencies from toolshed",
			},
			"install_resolver_dependencies": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Install resolver dependencies",
			},
			"tool_panel_section_id": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"new_tool_panel_section_label"},
				ForceNew:      true,
				Description:   "Tool panel section ID to list tool under",
			},
			"new_tool_panel_section_label": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "",
				ConflictsWith: []string{"tool_panel_section_id"},
				ForceNew:      true,
				Description:   "Label of tool panel section to create and list tool under",
			},
			"remove_from_disk": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Repository files from disk on uninstall",
			},
		},
		Description: "Tools are bundled and installed as repositories made available via [Galaxy Toolshed](https://toolshed.g2.bx.psu.edu/) deployments. This resource represents and manages an installed repository within a Galaxy instance.",
	}
}

func resourceRepositoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	toolShed := d.Get("tool_shed").(string)
	owner := d.Get("owner").(string)
	name := d.Get("name").(string)
	revision := d.Get("changeset_revision").(string)

	if repos, err := repositories.Install(ctx, g,
		toolShed,
		owner,
		name,
		revision,
		d.Get("install_tool_dependencies").(bool),
		d.Get("install_repository_dependencies").(bool),
		d.Get("install_resolver_dependencies").(bool),
		d.Get("tool_panel_section_id").(string),
		d.Get("new_tool_panel_section_label").(string),
		600, // 10 minute timeout
	); err == nil {
		if repos == nil || len(repos) == 0 {
			return diag.Errorf("Repository %v/%v/%v/%v already installed", toolShed, owner, name, revision)
		}
		var diags diag.Diagnostics
		if len(repos) > 1 {
			var ids []string
			for _, repo := range repos {
				ids = append(ids, repo.GetID())
			}
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Unexpected number of repositories created: %v", len(repos)),
				Detail:   fmt.Sprintf("Repository IDs: %v", ids),
			})
		}

		// Flatten tool_shed_status
		if err := d.Set("latest_installable_revision", repos[0].ToolShedStatus.LatestInstallableRevision); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("revision_update", repos[0].ToolShedStatus.RevisionUpdate); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("revision_upgrade", repos[0].ToolShedStatus.RevisionUpgrade); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("repository_deprecated", repos[0].ToolShedStatus.RepositoryDeprecated); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}

		return append(diags, toSchema(repos[0], d, repositoryOmitFields)...)
	} else {
		return diag.FromErr(err)
	}
}

func resourceRepositoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	if repo, err := repositories.Get(ctx, g, d.Id()); err == nil {
		return toSchema(repo, d, repositoryOmitFields)
	} else {
		return diag.FromErr(err)
	}
}

func resourceRepositoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//g := m.(*blend4go.GalaxyInstance)
	return nil

	// TODO if install_resolver_dependencies, install_tool_dependencies, install_repository_dependencies, do so, all other should force update
	// TODO look into if new_tool_panel_section_label or tool_panel_section_id can be updated
}

func resourceRepositoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	if err := repositories.UninstallID(ctx, g, d.Id(), d.Get("remove_from_disk").(bool)); err == nil {
		return nil
	} else {
		return diag.FromErr(err)
	}
}
