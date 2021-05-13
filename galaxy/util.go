package galaxy

import (
	"crypto/sha1"
	"encoding/base64"
	"github.com/brinkmanlab/blend4go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
	"strings"
)

var SetType = reflect.TypeOf((*schema.Set)(nil))

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

func fromSchema(m blend4go.GalaxyModel, s *schema.ResourceData, omit *map[string]interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	v := reflect.ValueOf(m)
	t := reflect.Indirect(v).Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if tag, ok := f.Tag.Lookup("json"); ok {
			if name := strings.Split(tag, ",")[0]; name != "" && name != "-" && name != "id" {
				value := reflect.ValueOf(s.Get(name))
				if omit != nil {
					if _, omitted := (*omit)[name]; omitted {
						continue
					}
				}
				if s.Get(name) == nil { // This is why we do not need the omit param as toSchema does
					continue
				}
				if f.Type.Kind() == reflect.Interface {
					continue // Fields without concrete type can't be handled
				}
				if value.Type().ConvertibleTo(f.Type) {
					reflect.Indirect(v).FieldByName(f.Name).Set(value.Convert(f.Type))
				} else {
					switch value.Kind() {
					case reflect.Map:
						newValue := reflect.MakeMapWithSize(f.Type, value.Len())
						iter := value.MapRange()
						for iter.Next() {
							newValue.SetMapIndex(iter.Key(), iter.Value().Elem().Convert(f.Type.Elem()))
						}
						reflect.Indirect(v).FieldByName(f.Name).Set(reflect.Indirect(newValue))
					case reflect.Ptr, reflect.Struct: // Pointer to struct?
						switch value.Type() {
						case SetType:
							value = value.MethodByName("List").Call(nil)[0]
						default:
							goto UNEXPECTED
						}
						fallthrough
					case reflect.Array, reflect.Slice:
						newValue := reflect.MakeSlice(f.Type, value.Len(), value.Len())
						for i := 0; i < value.Len(); i++ {
							newValue.Index(i).Set(value.Index(i).Elem().Convert(f.Type.Elem()))
						}
						reflect.Indirect(v).FieldByName(f.Name).Set(reflect.Indirect(newValue))
					default:
						goto UNEXPECTED
					}
					continue
				UNEXPECTED:
					diags = append(diags, diag.Errorf("Unexpected schema type (%v), from: %v, to: %v", name, value.Kind(), v.Kind())...)
				}
			}
		}
	}

	m.SetID(s.Id())
	return diags
}

func HashString(value string) string {
	h := sha1.New()
	h.Write([]byte(value))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
