package galaxy

import (
	"github.com/brinkmanlab/blend4go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
	"strings"
)

func toSchema(m blend4go.GalaxyModel, s *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	t := reflect.TypeOf(m)
	v := reflect.ValueOf(m)
	s.SetId(m.GetID())

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if tag, ok := f.Tag.Lookup("json"); ok {
			if name := strings.Split(tag, ",")[0]; name != "" && name != "-" && name != "id" {
				if err := s.Set(name, reflect.Indirect(v).FieldByName(f.Name)); err != nil {
					diags = append(diags, diag.FromErr(err)...)
				}
			}
		}
	}

	return diags
}

func fromSchema(m blend4go.GalaxyModel, s *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	t := reflect.TypeOf(m)
	v := reflect.ValueOf(m)
	m.SetID(s.Id())

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if tag, ok := f.Tag.Lookup("json"); ok {
			if name := strings.Split(tag, ",")[0]; name != "" && name != "-" && name != "id" {
				if err := s.Set(name, reflect.Indirect(v).FieldByName(f.Name)); err != nil {
					diags = append(diags, diag.FromErr(err)...)
				}
			}
		}
	}

	return diags
}
