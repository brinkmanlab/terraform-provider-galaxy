package galaxy

import (
	"context"
	"github.com/brinkmanlab/blend4go"
	"github.com/brinkmanlab/blend4go/users"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var userOmitFields = map[string]interface{}{"preferences": nil}

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			//"id": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username to identify user",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password to authenticate user against Galaxy",
			},
			"quota_percent": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Storage quota, between 0 and 100",
			},
			/*"preferences": &schema.Schema{
				Type:     ?,
				Computed: true,
			},*/
			"total_disk_usage": &schema.Schema{
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Total disk usage of users stored data",
			},
			"deleted": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "User deleted",
			},
			"purged": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "User purged",
			},
			"nice_total_disk_usage": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Human readable total disk usage of users stored data",
			},
			"quota": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Maximum disk storage available to user",
			},
			"email": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Users email address",
			},
			"is_admin": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "User is administrator",
			},
			"tags_used": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of tags assigned to users resources",
			},
			"api_key": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "API key of user",
			},
			"purge": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Purge a user on deletion",
			},
		},
		Description: "Create and manage Galaxy users. Used mostly to configure admin users.",
	}
}

func handleUser(ctx context.Context, user *users.User, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	if apiKey, err := user.GetAPIKey(ctx, d.Get("password").(string)); err == nil {
		if err := d.Set("api_key", apiKey); err != nil {
			diags = diag.FromErr(err)
		}
	} else {
		diags = diag.FromErr(err)
	}
	return append(diags, toSchema(user, d, userOmitFields)...)
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	if user, err := users.NewUser(ctx, g, d.Get("username").(string), d.Get("password").(string), d.Get("email").(string)); err == nil {
		return handleUser(ctx, user, d)
	} else {
		if err, ok := err.(*blend4go.ErrorResponse); ok {
			if err.Code == 400008 {
				// Attempt to undelete user
				if res, e := users.List(ctx, g, true, "", d.Get("username").(string), ""); e == nil {
					if len(res) == 0 {
						return diag.FromErr(err)
					}
					if len(res) != 1 {
						return diag.Errorf("unexpected number of results when searching for deleted user: %v", len(res))
					}
					if err := res[0].Undelete(ctx); err == nil {
						return handleUser(ctx, res[0], d)
					} else {
						return diag.FromErr(err)
					}
				} else {
					return diag.FromErr(e)
				}
			} else {
				return diag.FromErr(err)
			}
		}
		return diag.FromErr(err)
	}
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	if user, err := users.Get(ctx, g, d.Id(), false); err == nil {
		return toSchema(user, d, userOmitFields)
	} else {
		return diag.FromErr(err)
	}
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	g := m.(*blend4go.GalaxyInstance)

	user := new(users.User)
	user.SetGalaxyInstance(g)
	diags = append(diags, fromSchema(user, d)...)
	if d.HasChange("password") {
		current, new := d.GetChange("password")
		if err := user.SetPassword(ctx, current.(string), new.(string)); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}
	if err := user.Update(ctx); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	diags = append(diags, toSchema(user, d, userOmitFields)...)

	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	g := m.(*blend4go.GalaxyInstance)
	user := new(users.User)
	user.SetGalaxyInstance(g)
	diags = append(diags, fromSchema(user, d)...)
	if err := user.Delete(ctx, d.Get("purge").(bool)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
