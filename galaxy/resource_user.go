package galaxy

import (
	"context"
	"github.com/brinkmanlab/blend4go"
	"github.com/brinkmanlab/blend4go/users"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"quota_percent": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			/*"preferences": &schema.Schema{
				Type:     ?,
				Computed: true,
			},*/
			"total_disk_usage": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"deleted": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"purged": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"nice_total_disk_usage": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"quota": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"is_admin": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tags_used": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schema.TypeString,
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	if user, err := users.NewUser(ctx, g, d.Get("username").(string), d.Get("password").(string), d.Get("email").(string)); err == nil {
		return toSchema(user, d)
	} else {
		return diag.FromErr(err)
	}
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	if user, err := users.Get(ctx, g, d.Get("id").(string)); err == nil {
		return toSchema(user, d)
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
	diags = append(diags, diag.FromErr(user.Update(ctx))...)
	diags = append(diags, toSchema(user, d)...)

	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	g := m.(*blend4go.GalaxyInstance)
	user := new(users.User)
	user.SetGalaxyInstance(g)
	diags = append(diags, fromSchema(user, d)...)
	diags = append(diags, diag.FromErr(user.Delete(ctx))...)

	return diags
}
