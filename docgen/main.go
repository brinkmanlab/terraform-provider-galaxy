package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"terraform-provider-galaxy/galaxy"
	"text/template"
)

// Generate documentation from provider schemas
func main() {
	s := galaxy.Provider()
	tmpl := template.New("")
	type resource struct {
		Name                string
		Schema              map[string]*schema.Schema
		DepreciationMessage string
		Description         string
		Level               int
	}
	tmpl.Funcs(template.FuncMap{
		"inc": func(x int) int { return x + 1 },
		"listAttr": func(a []string, sep string) string {
			for i := range a {
				a[i] = "`" + a[i] + "`"
			}
			if len(a) > 1 {
				return strings.Join(a[:len(a)-1], ", ") + " " + sep + " " + a[len(a)-1]
			} else {
				return a[0]
			}
		},
		"typeName":   func(t schema.ValueType) string { return strings.TrimPrefix(t.String(), "Type") },
		"isResource": func(x interface{}) bool { _, ok := x.(*schema.Resource); return ok },
		"i":          func(x int) string { return strings.Repeat("  ", x) },
		"example": func(name string) (string, error) {
			if ex, err := ioutil.ReadFile(path.Join("./docgen/examples", name+".tf")); err == nil {
				return string(ex), nil
			} else {
				return "", err
			}
		},
		"tmplParams": func(name string, schema map[string]*schema.Schema, depMsg string, desc string, level int) *resource {
			return &resource{
				Name:                name,
				Schema:              schema,
				DepreciationMessage: depMsg,
				Description:         desc,
				Level:               level,
			}
		},
	})
	if tmpl, err := tmpl.ParseGlob("./docgen/*.md"); err == nil {
		if file, err := os.Create("./docs/index.md"); err == nil {
			if err := tmpl.ExecuteTemplate(file, "index.md", &resource{Name: "Galaxy", Schema: s.Schema, Description: galaxy.Description}); err != nil {
				panic(err)
			}
			if err := file.Close(); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
		for name, source := range s.DataSourcesMap {
			if file, err := os.Create("./docs/data-sources/" + name + ".md"); err == nil {
				if err := tmpl.ExecuteTemplate(file, "data_source.md", &resource{Name: name, Schema: source.Schema, DepreciationMessage: source.DeprecationMessage, Description: source.Description}); err != nil {
					panic(err)
				}
				if err := file.Close(); err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}
		}
		for name, source := range s.ResourcesMap {
			if file, err := os.Create("./docs/resources/" + name + ".md"); err == nil {
				if err := tmpl.ExecuteTemplate(file, "resource.md", &resource{Name: name, Schema: source.Schema, DepreciationMessage: source.DeprecationMessage, Description: source.Description}); err != nil {
					panic(err)
				}
				if err := file.Close(); err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}
		}
	} else {
		panic(err)
	}
}
