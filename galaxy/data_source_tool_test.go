package galaxy_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

const ToolPath = "./test-fixtures/tool.tf"

func TestAccTool_basic(t *testing.T) {
	tmpl := testAccConfigTemplate(ToolPath, t)
	name := "test"
	resourceName := "data.galaxy_tool." + name
	type tmplFields struct {
		Name string
		Id   string
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(tmpl, t, &tmplFields{Name: name, Id: "toolshed.g2.bx.psu.edu/repos/brinkmanlab/awkscript/awkscript/1.0"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "AWK Script"),
				),
			},
		},
	})
}
