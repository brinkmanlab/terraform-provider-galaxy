package galaxy

import (
	"context"
	"github.com/brinkmanlab/blend4go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"net/http"
	"time"
)

// TODO If schema.Provider.Description is ever added, move this there
// https://github.com/hashicorp/terraform-plugin-sdk/issues/631
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
			"apikey": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				DefaultFunc:  schema.EnvDefaultFunc("GALAXY_API_KEY", nil),
				ExactlyOneOf: []string{"apikey", "username"},
				Description:  "API key associated with a Galaxy administrator account. A master API key will fail to create resources that need to be associated with a user. Refers to GALAXY_API_KEY env variable if unset.",
			},
			"username": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("GALAXY_USERNAME", nil),
				ExactlyOneOf: []string{"apikey", "username"},
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
			"wait_for_host": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GALAXY_WAIT", nil),
				Description: "Some terraform resources return prematurely causing this provider to fail to resolve the Galaxy host. Specify in seconds how long to wait for the host to become available (0 for forever)",
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
	if host, ok := d.GetOk("host"); ok && host.(string) != "" {
		if wait, ok := d.GetOkExists("wait_for_host"); ok {
			duration := time.Duration(wait.(int)) * time.Second
			for start := time.Now(); (wait == 0 || time.Since(start) < duration) && ctx.Err() == nil; time.Sleep(time.Second * 2) {
				if res, err := http.Get(host.(string)); err == nil {
					if res.StatusCode < 400 {
						break
					}
				}
				log.Printf("[INFO] Waiting for Galaxy host..")
			}
		}
		key, ok := d.GetOk("apikey")
		if !ok {
			if user, ok := d.GetOk("username"); ok {
				var err error
				if key, err = blend4go.GetAPIKey(ctx, host.(string), user.(string), d.Get("password").(string)); err != nil {
					return nil, diag.FromErr(err)
				}
			} else {
				return nil, diag.Errorf("API key or username must be provided and non-empty")
			}
		}

		var level blend4go.LogLevel
		switch logging.LogLevel() {
		case "":
		default:
			level = blend4go.NONE
			break
		case "TRACE":
			level = blend4go.DEBUG
			break
		case "DEBUG":
			level = blend4go.DEBUG
			break
		case "INFO":
			level = blend4go.INFO
			break
		case "WARN":
			level = blend4go.WARN
			break
		case "ERROR":
			level = blend4go.ERROR
			break
		}

		c := blend4go.NewGalaxyInstanceLogger(host.(string), key.(string), log.Writer(), level)

		// Test connection
		if _, err := c.Version(ctx); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}

		return c, diags
	} else {
		return nil, diag.Errorf("Galaxy host URL must be provided and non-empty")
	}
}
