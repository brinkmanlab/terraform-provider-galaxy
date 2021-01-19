package galaxy

import (
	"context"
	"github.com/brinkmanlab/blend4go"
	"github.com/brinkmanlab/blend4go/repositories"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
)

var repositoryOmitFields = map[string]interface{}{"tool_shed_status": nil}

func resourceRepository() *schema.Resource {
	subRepo := map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Galaxy ID of installed repository",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Installation status",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Repository name",
			ForceNew:    true,
		},
		"deleted": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Repository deleted",
		},
		"ctx_rev": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"error_message": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Install error message",
		},
		"installed_changeset_revision": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Installed changeset revision",
		},
		"tool_shed": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Repository toolshed",
		},
		"dist_to_shed": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"url": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Repository url",
		},
		"uninstalled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Uninstalled",
		},
		"owner": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Repository owner",
		},
		"changeset_revision": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Changeset revision of repository",
		},
		"include_datatypes": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Repository includes datatypes",
		},
		"latest_installable_revision": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Latest installable revision of repository",
		},
		"revision_update": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "https://github.com/galaxyproject/galaxy/issues/10453",
		},
		"revision_upgrade": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "https://github.com/galaxyproject/galaxy/issues/10453",
		},
		"repository_deprecated": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Repository depreciated. https://github.com/galaxyproject/galaxy/issues/10453",
		},
		"tools": {
			Type:        schema.TypeSet, // TODO change to map[tool_id] after https://github.com/hashicorp/terraform-plugin-sdk/issues/633
			Computed:    true,
			Description: "List of tools installed by repository",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"tool_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Tool Id",
					},
					"tool_guid": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Tool guid",
					},
					"name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Tool name",
					},
					"version": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Tool version",
					},
					"description": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Tool description",
					},
					"config_file": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Path to tool wrapper XML (on toolshed)",
					},
				},
			},
		},
	}

	repo := map[string]*schema.Schema{
		"install_tool_dependencies": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Install tool dependencies using the configured dependency manager",
			ForceNew:    true,
		},
		"install_repository_dependencies": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
			//Deprecated:  "Repository dependencies are depreciated",
			Description: "Install repository dependencies from toolshed",
			ForceNew:    true,
		},
		"install_resolver_dependencies": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Install resolver dependencies",
			ForceNew:    true,
		},
		"tool_panel_section_id": {
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"new_tool_panel_section_label"},
			ForceNew:      true,
			Description:   "Tool panel section ID to list tool under",
		},
		"new_tool_panel_section_label": {
			Type:          schema.TypeString,
			Optional:      true,
			Default:       "",
			ConflictsWith: []string{"tool_panel_section_id"},
			ForceNew:      true,
			Description:   "Label of tool panel section to create and list tool under",
		},
		"remove_from_disk": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Repository files from disk on uninstall",
			ForceNew:    true,
		},
		"sub_repositories": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				ReadContext:   resourceJobRead,
				DeleteContext: resourceJobDelete,
				Schema:        subRepo,
			},
			Description: "If the repository is a suite or otherwise depends on other repositories, they will be listed here",
		},
	}

	// Copy in common repo fields
	for k, v := range subRepo {
		if _, ok := repo[k]; k != "id" && !ok { // Do not overwrite redeclared fields
			repo[k] = v
		}
	}

	return &schema.Resource{
		CreateContext: resourceRepositoryCreate,
		ReadContext:   resourceRepositoryRead,
		//UpdateContext: resourceRepositoryUpdate,
		DeleteContext: resourceRepositoryDelete,
		Schema:        repo,
		Description:   "Tools are bundled and installed as repositories made available via [Galaxy Toolshed](https://toolshed.g2.bx.psu.edu/) deployments. This resource represents and manages an installed repository within a Galaxy instance.",
	}
}

// Populate tools
func populateTools(ctx context.Context, d *schema.ResourceData, repo *repositories.Repository) diag.Diagnostics {
	if tools, err := repo.Tools(ctx); err == nil {
		r := make([]map[string]string, len(tools))
		for i, tool := range tools {
			r[i] = map[string]string{
				"tool_id":     tool.Id,
				"tool_guid":   tool.Guid,
				"name":        tool.Name,
				"version":     tool.Version,
				"description": tool.Description,
				"config_file": tool.ConfigFile,
			}
		}
		if err := d.Set("tools", r); err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.FromErr(err)
	}
	return nil
}

func flattenRepository(d *schema.ResourceData, repo *repositories.Repository) diag.Diagnostics {
	var diags diag.Diagnostics
	// Flatten tool_shed_status
	if err := d.Set("latest_installable_revision", repo.ToolShedStatus.LatestInstallableRevision); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("revision_update", repo.ToolShedStatus.RevisionUpdate); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("revision_upgrade", repo.ToolShedStatus.RevisionUpgrade); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("repository_deprecated", repo.ToolShedStatus.RepositoryDeprecated); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	return diags
}

func repoToMap(repo *repositories.Repository) (map[string]interface{}, error) {
	sr := map[string]interface{}{}
	if decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{TagName: "json", Result: &sr}); err == nil {
		if err := decoder.Decode(repo); err != nil {
			return nil, err
		}
		delete(sr, "tool_shed_status")
		sr["latest_installable_revision"] = repo.ToolShedStatus.LatestInstallableRevision
		sr["revision_update"] = repo.ToolShedStatus.RevisionUpdate
		sr["revision_upgrade"] = repo.ToolShedStatus.RevisionUpgrade
		sr["repository_deprecated"] = repo.ToolShedStatus.RepositoryDeprecated
		return sr, nil
	} else {
		return nil, err
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
		var subRepos []map[string]interface{}
		for _, repo := range repos {
			if repo.ToolShed == toolShed && repo.Owner == owner && repo.Name == name {
				diags = append(diags, flattenRepository(d, repo)...)
				diags = append(diags, populateTools(ctx, d, repo)...)
				diags = append(diags, toSchema(repo, d, repositoryOmitFields)...)
			} else {
				if subRepo, err := repoToMap(repo); err == nil {
					subRepos = append(subRepos, subRepo)
				} else {
					diags = append(diags, diag.FromErr(err)...)
				}
			}
		}

		if err := d.Set("sub_repositories", subRepos); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}

		return diags
	} else {
		return diag.FromErr(err)
	}
}

func resourceRepositoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	if repo, err := repositories.Get(ctx, g, d.Id()); err == nil {
		var diags diag.Diagnostics
		//diags = append(diags, populateTools(ctx, d, repo)...)
		diags = append(diags, toSchema(repo, d, repositoryOmitFields)...)
		if repos := d.Get("sub_repositories"); len(repos.([]interface{})) > 0 {
			var subRepos []map[string]interface{}
			for _, subRepo := range repos.([]interface{}) {
				if repo, err := repositories.Get(ctx, g, subRepo.(map[string]interface{})["id"].(string)); err == nil {
					if sr, err := repoToMap(repo); err == nil {
						subRepos = append(subRepos, sr)
					} else {
						diags = append(diags, diag.FromErr(err)...)
					}
				}
			}
			if err := d.Set("sub_repositories", subRepos); err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

		return diags
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
	var diags diag.Diagnostics

	if err := repositories.UninstallID(ctx, g, d.Id(), d.Get("remove_from_disk").(bool)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if repos := d.Get("sub_repositories"); len(repos.([]interface{})) > 0 {
		for _, subRepo := range repos.([]interface{}) {
			if err := repositories.UninstallID(ctx, g, subRepo.(map[string]interface{})["id"].(string), d.Get("remove_from_disk").(bool)); err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
		}
	}

	return diags
}
