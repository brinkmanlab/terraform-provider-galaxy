package galaxy_test

import (
	"context"
	"fmt"
	"github.com/brinkmanlab/blend4go/users"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

const UserResourcePath = "test-fixtures/user.tf"

func testAccUserExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ID unset")
		}

		if res, err := users.Get(context.Background(), testAccGalaxyInstance(), rs.Primary.ID, false); err == nil {
			if res.GetID() != rs.Primary.ID {
				return fmt.Errorf("ID mismatch between stored ID (%v) and fetched (%v)", rs.Primary.ID, res.GetID())
			}
		} else {
			return err
		}

		return nil
	}
}

func TestAccUser_basic(t *testing.T) {
	tmpl := testAccConfigTemplate(UserResourcePath, t)
	name := "test"
	resourceName := "galaxy_user." + name
	type tmplFields struct {
		Name     string
		Username string
		Password string
		Email    string
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		IDRefreshName:     resourceName,
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(tmpl, t, &tmplFields{Name: name, Username: "test", Password: "testpass", Email: "test@example.com"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "username", "test"),
					resource.TestCheckResourceAttr(resourceName, "email", "test@example.com"),
					resource.TestCheckResourceAttrSet(resourceName, "api_key"),
					//testCheckResourceAttrEqual(resourceName, "deleted", false),
					//testCheckResourceAttrEqual(resourceName, "purged", false),
				),
			},
		},
	})
}
