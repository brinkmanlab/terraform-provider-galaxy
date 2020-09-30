package galaxy

import (
	"context"
	"github.com/brinkmanlab/blend4go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GALAXY_HOST", nil),
				Description: "URL to Galaxy instance",
			},
			"api_key": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				DefaultFunc:  schema.EnvDefaultFunc("GALAXY_API_KEY", nil),
				ExactlyOneOf: []string{"api_key", "username"},
				Description:  "API key associated with a Galaxy administrator account. A master API key will fail to create resources that need to be associated with a user.",
			},
			"username": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("GALAXY_USERNAME", nil),
				ExactlyOneOf: []string{"api_key", "username"},
				RequiredWith: []string{"password"},
				Description:  "Username or email address of Galaxy administrator account",
			},
			"password": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				DefaultFunc:  schema.EnvDefaultFunc("GALAXY_PASSWORD", nil),
				RequiredWith: []string{"username"},
				Description:  "Password associated with username",
			},
		},
		ConfigureContextFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"galaxy_user":            resourceUser(),
			"galaxy_stored_workflow": resourceStoredWorkflow(),
			"galaxy_job":             resourceJob(),
			"galaxy_repository":      resourceRepository(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"galaxy_workflow_tools": dataSourceWorkflowRepositories(),
			"galaxy_tool":           dataSourceTool(),
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
