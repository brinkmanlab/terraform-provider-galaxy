package galaxy_test

import (
	"context"
	"fmt"
	"github.com/brinkmanlab/blend4go/jobs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

const JobResourcePath = "test-fixtures/job.tf"

func testAccJobExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ID unset")
		}

		if res, err := jobs.Get(context.Background(), testAccGalaxyInstance(), rs.Primary.ID); err == nil {
			if res.GetID() != rs.Primary.ID {
				return fmt.Errorf("ID mismatch between stored ID (%v) and fetched (%v)", rs.Primary.ID, res.GetID())
			}
		} else {
			return err
		}

		return nil
	}
}

func TestAccJob_basic(t *testing.T) {
	tmpl := testAccConfigTemplate(JobResourcePath, t)
	name := "test"
	resourceName := "galaxy_job." + name
	type tmplFields struct {
		Name    string
		JobName string
		ToolID  string
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(tmpl, t, &tmplFields{Name: name, JobName: "test", ToolID: "toolshed.g2.bx.psu.edu/repos/brinkmanlab/awkscript/awkscript/1.0"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccJobExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tool_id", "toolshed.g2.bx.psu.edu/repos/brinkmanlab/awkscript/awkscript/1.0"),
					resource.TestCheckResourceAttrSet(resourceName, "history_id"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					//testCheckResourceAttrEqual(resourceName, "deleted", false),
					//testCheckResourceAttrEqual(resourceName, "purged", false),
				),
			},
		},
	})
}
