package galaxy

import (
	"github.com/brinkmanlab/blend4go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
	"strings"
)

func toSchema(m blend4go.GalaxyModel, s *schema.ResourceData, omit map[string]interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	v := reflect.ValueOf(m)
	t := reflect.Indirect(v).Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if tag, ok := f.Tag.Lookup("json"); ok {
			if name := strings.Split(tag, ",")[0]; name != "" && name != "-" && name != "id" {
				if _, omitted := omit[name]; !omitted {
					if err := s.Set(name, reflect.Indirect(v).FieldByName(f.Name).Interface()); err != nil {
						diags = append(diags, diag.FromErr(err)...)
					}
				}
			}
		}
	}

	s.SetId(m.GetID())
	return diags
}

func fromSchema(m blend4go.GalaxyModel, s *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	v := reflect.ValueOf(m)
	t := reflect.Indirect(v).Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if tag, ok := f.Tag.Lookup("json"); ok {
			if name := strings.Split(tag, ",")[0]; name != "" && name != "-" && name != "id" {
				value := reflect.ValueOf(s.Get(name))
				if s.Get(name) == nil {
					continue
				}
				if f.Type.Kind() == reflect.Interface {
					continue // Fields without concrete type can't be handled
				}
				if (value.Kind() != reflect.Map && value.Kind() != reflect.Array && value.Kind() != reflect.Slice) || value.Type().Elem() == f.Type.Elem() {
					reflect.Indirect(v).FieldByName(f.Name).Set(value.Convert(f.Type))
				} else if value.Kind() == reflect.Map {
					new_value := reflect.MakeMapWithSize(f.Type, value.Len())
					iter := value.MapRange()
					for iter.Next() {
						new_value.SetMapIndex(iter.Key(), iter.Value().Elem().Convert(f.Type.Elem()))
					}
					reflect.Indirect(v).FieldByName(f.Name).Set(reflect.Indirect(new_value))
				} else if value.Kind() == reflect.Array || value.Kind() == reflect.Slice {
					new_value := reflect.MakeSlice(f.Type, value.Len(), value.Len())
					for i := 0; i < value.Len(); i++ {
						new_value.Index(i).Set(value.Index(i).Convert(f.Type.Elem()))
					}
				} else {
					diags = append(diags, diag.Errorf("Unexpected schema type (%v), from: %v, to: %v", name, value.Kind(), v.Kind())...)
				}
			}
		}
	}

	m.SetID(s.Id())
	return diags
}
