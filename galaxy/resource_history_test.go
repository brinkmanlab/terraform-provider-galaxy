package galaxy_test

import (
	"context"
	"fmt"
	"github.com/brinkmanlab/blend4go/histories"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

const HistoryResourcePath = "test-fixtures/history.tf"

func testAccHistoryExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ID unset")
		}

		if res, err := histories.Get(context.Background(), testAccGalaxyInstance(), rs.Primary.ID); err == nil {
			if res.GetID() != rs.Primary.ID {
				return fmt.Errorf("ID mismatch between stored ID (%v) and fetched (%v)", rs.Primary.ID, res.GetID())
			}
		} else {
			return err
		}

		return nil
	}
}

func TestAccHistory_basic(t *testing.T) {
	tmpl := testAccConfigTemplate(HistoryResourcePath, t)
	name := "test"
	resourceName := "galaxy_history." + name
	type tmplFields struct {
		name        string
		historyName string
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		IDRefreshName:     resourceName,
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(tmpl, t, &tmplFields{name: name, historyName: "test"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccHistoryExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test"),
					//testCheckResourceAttrEqual(resourceName, "deleted", false),
					//testCheckResourceAttrEqual(resourceName, "purged", false),
				),
			},
		},
	})
}
