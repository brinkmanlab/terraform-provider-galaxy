package galaxy

import (
	"context"
	"github.com/brinkmanlab/blend4go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO If schema.Provider.Description is ever added, move this there
const Description = `Manages resources of a [Galaxy](https://galaxyproject.org) instance.

Based on the [blend4go](https://github.com/brinkmanlab/blend4go) library for Galaxy API requests

Written and maintained by the [Fiona Brinkman Laboratory](https://github.com/brinkmanlab/terraform-provider-galaxy)
`

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GALAXY_HOST", nil),
				Description: "URL to Galaxy instance. Refers to GALAXY_HOST env variable if unset.",
			},
			"api_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				DefaultFunc:  schema.EnvDefaultFunc("GALAXY_API_KEY", nil),
				ExactlyOneOf: []string{"api_key", "username"},
				Description:  "API key associated with a Galaxy administrator account. A master API key will fail to create resources that need to be associated with a user. Refers to GALAXY_API_KEY env variable if unset.",
			},
			"username": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("GALAXY_USERNAME", nil),
				ExactlyOneOf: []string{"api_key", "username"},
				RequiredWith: []string{"password"},
				Description:  "Username or email address of Galaxy administrator account. Refers to GALAXY_USERNAME env variable if unset.",
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				DefaultFunc:  schema.EnvDefaultFunc("GALAXY_PASSWORD", nil),
				RequiredWith: []string{"username"},
				Description:  "Password associated with username. Refers to GALAXY_PASSWORD env variable if unset.",
			},
		},
		ConfigureContextFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"galaxy_user":            resourceUser(),
			"galaxy_stored_workflow": resourceStoredWorkflow(),
			"galaxy_job":             resourceJob(),
			"galaxy_repository":      resourceRepository(),
			"galaxy_history":         resourceHistory(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"galaxy_workflow_repositories": dataSourceWorkflowRepositories(),
			"galaxy_tool":                  dataSourceTool(),
		},
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	host := d.Get("host").(string)
	key, ok := d.GetOk("api_key")
	if !ok {
		var err error
		key, err = blend4go.GetAPIKey(ctx, host, d.Get("username").(string), d.Get("password").(string))
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	c := blend4go.NewGalaxyInstance(host, key.(string))

	return c, diags
}
