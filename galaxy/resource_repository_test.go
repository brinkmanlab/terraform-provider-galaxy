package galaxy_test

import (
	"context"
	"fmt"
	"github.com/brinkmanlab/blend4go/repositories"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

const RepositoryResourcePath = "test-fixtures/repository.tf"

func testAccRepositoryExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ID unset")
		}

		if res, err := repositories.Get(context.Background(), testAccGalaxyInstance(), rs.Primary.ID); err == nil {
			if res.GetID() != rs.Primary.ID {
				return fmt.Errorf("ID mismatch between stored ID (%v) and fetched (%v)", rs.Primary.ID, res.GetID())
			}
		} else {
			return err
		}

		return nil
	}
}

func TestAccRepository_basic(t *testing.T) {
	tmpl := testAccConfigTemplate(RepositoryResourcePath, t)
	name := "test"
	resourceName := "galaxy_repository." + name
	type tmplFields struct {
		Name              string
		RepositoryName    string
		Toolshed          string
		Owner             string
		RepoName          string
		ChangesetRevision string
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		IDRefreshName:     resourceName,
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(tmpl, t, &tmplFields{Name: name, RepositoryName: "test", Toolshed: "toolshed.g2.bx.psu.edu", Owner: "brinkmanlab", RepoName: "awkscript", ChangesetRevision: "7966a43dbc9e"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccRepositoryExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tool_shed", "toolshed.g2.bx.psu.edu"),
					resource.TestCheckResourceAttr(resourceName, "owner", "brinkmanlab"),
					resource.TestCheckResourceAttr(resourceName, "name", "awkscript"),
					//resource.TestCheckResourceAttr(resourceName, "changeset_revision", "7966a43dbc9e"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					testCheckResourceAttrEqual(resourceName, "tools.#", 1),
					resource.TestCheckResourceAttr(resourceName, "tools.0.tool_id", "awkscript"),
					//testCheckResourceAttrEqual(resourceName, "deleted", false),
				),
			},
		},
	})
}
