package galaxy

import (
	"context"
	"github.com/brinkmanlab/blend4go"
	"github.com/brinkmanlab/blend4go/groups"
	"github.com/brinkmanlab/blend4go/quotas"
	"github.com/brinkmanlab/blend4go/users"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

var quotaOmitFields = map[string]interface{}{"default": nil, "users": nil, "groups": nil}

func resourceQuota() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceQuotaCreate,
		ReadContext:   resourceQuotaRead,
		UpdateContext: resourceQuotaUpdate,
		DeleteContext: resourceQuotaDelete,
		Schema: map[string]*schema.Schema{
			//"id": {
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Quota name as displayed to user",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Description of quota",
			},
			"operation": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(quotas.SetTo),
				ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
					switch v.(string) {
					default:
						diags := diag.Errorf("invalid operation %s", v)
						diags[0].AttributePath = path
						return diags
					case string(quotas.SetTo):
					case string(quotas.IncreaseBy):
					case string(quotas.DecreaseBy):
					}
					return nil
				},
				Description: "Assign (=), increase by amount (+), or decrease by amount (-)",
			},
			"amount": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Examples: \"10000MB\", \"99 gb\", \"0.2T\", \"unlimited\"",
			},
			"display_amount": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Human readable amount",
			},
			"bytes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Amount, in bytes",
			},
			"default": {
				Type:     schema.TypeString,
				Default:  string(quotas.NotDefault),
				Optional: true,
				ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
					switch v.(string) {
					default:
						diags := diag.Errorf("invalid default users value %s", v)
						diags[0].AttributePath = path
						return diags
					case string(quotas.NotDefault):
					case string(quotas.RegisteredUsers):
					case string(quotas.UnregisteredUsers):
					}
					return nil
				},
				Description: "Set as default for category of users (unregistered, registered)",
			},
			//"deleted": {
			//	Type:        schema.TypeBool,
			//	Computed:    true,
			//	Description: "Deleted",
			//},
			"users": {
				Type:         schema.TypeList,
				Elem:         &schema.Schema{Type: schema.TypeString},
				Optional:     true,
				AtLeastOneOf: []string{"users", "groups", "default"},
				Description:  "List of user ids to apply quota to",
			},
			"groups": {
				Type:         schema.TypeList,
				Elem:         &schema.Schema{Type: schema.TypeString},
				Optional:     true,
				AtLeastOneOf: []string{"users", "groups", "default"},
				Description:  "List of group ids to apply quota to",
			},
			"purge": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Purge a user on deletion",
			},
		},
		Importer:    &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Description: "Galaxy quotas regulate the amount of data a user can store in their account at any one time.",
	}
}

func quotaFromSchema(quota *quotas.Quota, d *schema.ResourceData) diag.Diagnostics {
	diags := fromSchema(quota, d, &quotaOmitFields)
	quota.Default = quotas.NotDefault
	switch d.Get("default").(string) {
	default:
		return diag.Errorf("invalid default users value %s", d.Get("default").(string))
	case string(quotas.NotDefault):
	case string(quotas.RegisteredUsers):
		quota.Default = quotas.RegisteredUsers
	case string(quotas.UnregisteredUsers):
		quota.Default = quotas.UnregisteredUsers
	}
	quota.Users = []*users.User{}
	for _, userID := range d.Get("users").([]interface{}) {
		quota.Users = append(quota.Users, &users.User{Id: userID.(blend4go.GalaxyID)})
	}
	quota.Groups = []*groups.Group{}
	for _, groupID := range d.Get("groups").([]interface{}) {
		quota.Groups = append(quota.Groups, &groups.Group{Id: groupID.(blend4go.GalaxyID)})
	}
	return diags
}

func quotaToSchema(quota *quotas.Quota, d *schema.ResourceData) diag.Diagnostics {
	diags := toSchema(quota, d, quotaOmitFields)
	if err := d.Set("default", string(quota.Default)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	var userIDs []blend4go.GalaxyID
	for _, user := range quota.Users {
		userIDs = append(userIDs, user.GetID())
	}
	if err := d.Set("users", userIDs); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	var groupIDs []blend4go.GalaxyID
	for _, group := range quota.Groups {
		groupIDs = append(groupIDs, group.GetID())
	}
	if err := d.Set("groups", groupIDs); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	return diags
}

func resourceQuotaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)
	operation := quotas.SetTo
	switch d.Get("operation").(string) {
	default:
		return diag.Errorf("invalid operation %s", operation)
	case string(quotas.SetTo):
	case string(quotas.IncreaseBy):
		operation = quotas.IncreaseBy
	case string(quotas.DecreaseBy):
		operation = quotas.DecreaseBy
	}

	default_for := quotas.NotDefault
	switch d.Get("default").(string) {
	default:
		return diag.Errorf("invalid default users value %s", default_for)
	case string(quotas.NotDefault):
	case string(quotas.RegisteredUsers):
		default_for = quotas.RegisteredUsers
	case string(quotas.UnregisteredUsers):
		default_for = quotas.UnregisteredUsers
	}

	var userIDs []blend4go.GalaxyID
	for _, user := range d.Get("users").([]interface{}) {
		userIDs = append(userIDs, user.(blend4go.GalaxyID))
	}
	var groupIDs []blend4go.GalaxyID
	for _, group := range d.Get("groups").([]interface{}) {
		groupIDs = append(groupIDs, group.(blend4go.GalaxyID))
	}

	name := d.Get("name").(string)
	amount := d.Get("amount").(string)
	description := d.Get("description").(string)
	if quota, err := quotas.NewQuota(ctx, g, name, amount, description, operation, userIDs, groupIDs, default_for); err == nil {
		return quotaToSchema(quota, d)
	} else {
		if strings.Contains(err.Error(), "quota with that name already exists") { // TODO https://github.com/galaxyproject/galaxy/issues/11971
			// Attempt to undelete quota
			if res, e := quotas.List(ctx, g, true); e == nil {
				for _, quota := range res {
					if quota.Name == name {
						if err := quota.Undelete(ctx); err == nil {
							var diags diag.Diagnostics
							quota.Name = name
							quota.Operation = operation
							quota.Description = description
							quota.Users = []*users.User{}
							for _, userID := range userIDs {
								quota.Users = append(quota.Users, &users.User{Id: userID})
							}
							quota.Groups = []*groups.Group{}
							for _, groupID := range groupIDs {
								quota.Groups = append(quota.Groups, &groups.Group{Id: groupID})
							}
							quota.Default = default_for
							if err := quota.Update(ctx, amount); err != nil {
								diags = append(diags, diag.FromErr(err)...)
							}
							diags = append(diags, quotaToSchema(quota, d)...)
							return diags
						} else {
							return diag.FromErr(err)
						}
					}
				}
				return diag.FromErr(err)
			} else {
				return diag.FromErr(e)
			}
		}
		return diag.FromErr(err)
	}
}

func resourceQuotaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	if quota, err := quotas.Get(ctx, g, d.Id(), false); err == nil {
		return quotaToSchema(quota, d)
	} else {
		return diag.FromErr(err)
	}
}

func resourceQuotaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)

	quota := new(quotas.Quota)
	quota.SetGalaxyInstance(g)
	diags := quotaFromSchema(quota, d)
	if err := quota.Update(ctx, ""); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	diags = append(diags, quotaToSchema(quota, d)...)
	return diags
}

func resourceQuotaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	g := m.(*blend4go.GalaxyInstance)
	quota := new(quotas.Quota)
	quota.SetGalaxyInstance(g)
	diags := quotaFromSchema(quota, d)
	if err := quota.Delete(ctx, d.Get("purge").(bool)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
