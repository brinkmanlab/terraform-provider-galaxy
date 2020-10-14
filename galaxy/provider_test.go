package galaxy_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/brinkmanlab/blend4go"
	"github.com/brinkmanlab/blend4go/test_util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"html/template"
	"os"
	"strings"
	"terraform-provider-galaxy/galaxy"
	"testing"
)

var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"galaxy": func() (*schema.Provider, error) {
			provider := galaxy.Provider()
			err := provider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))
			for _, e := range err {
				if e.Severity == diag.Error {
					return provider, errors.New(strings.Join([]string{e.Summary, e.Detail}, "\n"))
				}
			}
			return provider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := galaxy.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = galaxy.Provider()
}

func testAccPreCheck(t *testing.T) func() {
	return func() {
		if os.Getenv("GALAXY_HOST") == "" {
			t.Fatal("GALAXY_HOST must be set for acceptance tests")
		}

		if os.Getenv("GALAXY_API_KEY") == "" {
			t.Fatal("GALAXY_API_KEY must be set for acceptance tests")
		}

		provider, err := testAccProviderFactories["galaxy"]()
		if err != nil {
			t.Fatal(err)
			return
		}

		diags := provider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))
		if diags != nil && diags.HasError() {
			t.Fatal(diags)
		}
	}
}

func testAccConfigTemplate(file string, t *testing.T) *template.Template {
	if tmpl, err := template.ParseFiles(file); err == nil {
		return tmpl
	} else {
		t.Fatalf("Failed to parse config template: %v", file)
	}
	return nil
}

func testAccConfig(tmpl *template.Template, t *testing.T, data interface{}) string {
	var res strings.Builder
	if err := tmpl.Execute(&res, data); err == nil {
		return res.String()
	} else {
		t.Fatalf("Failed to populate config template %v with data: %v", t.Name(), data)
	}
	return ""
}

func testAccGalaxyInstance() *blend4go.GalaxyInstance {
	return test_util.NewTestInstance()
}

func testCheckResourceAttrEqual(resourceName, key string, value interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_ /*rs*/, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		return nil
	}
}
