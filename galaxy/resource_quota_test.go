package galaxy_test

import (
	"context"
	"fmt"
	"github.com/brinkmanlab/blend4go/quotas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

const QuotaResourcePath = "test-fixtures/quota.tf"
const QuotaUsersResourcePath = "test-fixtures/quota_users.tf"

func testAccQuotaExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ID unset")
		}

		if res, err := quotas.Get(context.Background(), testAccGalaxyInstance(), rs.Primary.ID, false); err == nil {
			if res.GetID() != rs.Primary.ID {
				return fmt.Errorf("ID mismatch between stored ID (%v) and fetched (%v)", rs.Primary.ID, res.GetID())
			}
		} else {
			return err
		}

		return nil
	}
}

func TestAccQuota_basic(t *testing.T) {
	tmpl := testAccConfigTemplate(QuotaResourcePath, t)
	name := "test"
	resourceName := "galaxy_quota." + name
	type tmplFields struct {
		Name       string
		Quotaname  string
		Amount     string
		DefaultFor string
	}
	/*resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		IDRefreshName:     resourceName,
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(tmpl, t, &tmplFields{Name: name, Quotaname: "test", Amount: "1000", DefaultFor: "registered"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccQuotaExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test"),
					resource.TestCheckResourceAttr(resourceName, "default", "registered"),
					testCheckResourceAttrEqual(resourceName, "bytes", 1000),
				),
			},
		},
	})*/
	tmpl = testAccConfigTemplate(QuotaUsersResourcePath, t)
	type userTmplFields struct {
		tmplFields
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
				Config: testAccConfig(tmpl, t, &userTmplFields{
					tmplFields: tmplFields{Name: name, Quotaname: "test", Amount: "1000", DefaultFor: "no"},
					Username:   "quota_test3",
					Password:   "testpassword",
					Email:      "quota_test3@test.com",
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccQuotaExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test"),
					resource.TestCheckResourceAttr(resourceName, "default", "no"),
					testCheckResourceAttrEqual(resourceName, "bytes", 1000),
					testCheckResourceAttrEqual(resourceName, "users.#", 1),
				),
			},
		},
	})
}
