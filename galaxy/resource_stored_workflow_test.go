package galaxy_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/brinkmanlab/blend4go/workflows"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"io/ioutil"
	"testing"
)

const WorkflowPath = "./test-fixtures/workflow.ga"
const WorkflowResourcePath = "./test-fixtures/workflow.tf"

func testAccWorkflowExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ID unset")
		}

		if res, err := workflows.Get(context.Background(), testAccGalaxyInstance(), rs.Primary.ID); err == nil {
			if res.GetID() != rs.Primary.ID {
				return fmt.Errorf("ID mismatch between stored ID (%v) and fetched (%v)", rs.Primary.ID, res.GetID())
			}
		} else {
			return err
		}

		return nil
	}
}

func loadWorkflow(path string) (string, map[string]interface{}, error) {
	if workflow, err := ioutil.ReadFile(path); err == nil {
		parsedWorkflow := make(map[string]interface{})
		if err := json.Unmarshal(workflow, &parsedWorkflow); err != nil {
			return "", nil, err
		}
		return string(workflow), parsedWorkflow, nil
	} else {
		return "", nil, err
	}
}

func TestAccWorkflow_basic(t *testing.T) {
	tmpl := testAccConfigTemplate(WorkflowResourcePath, t)
	name := "test"
	resourceName := "galaxy_stored_workflow." + name
	workflow, parsedWorkflow, err := loadWorkflow(WorkflowPath)
	if err != nil {
		t.Fatal(err)
	}

	type tmplFields struct {
		Name string
		Json string
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		IDRefreshName:     resourceName,
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(tmpl, t, &tmplFields{Name: name, Json: workflow}),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccWorkflowExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "json", workflow),
					resource.TestCheckResourceAttr(resourceName, "name", parsedWorkflow["name"].(string)),
					resource.TestCheckResourceAttr(resourceName, "annotation", parsedWorkflow["annotation"].(string)),
					//testCheckResourceAttrEquals(resourceName, "tags", parsed_workflow["tags"].([]string)),
				),
			},
		},
	})
}
