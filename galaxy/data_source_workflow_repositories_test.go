package galaxy_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"strings"
	"testing"
)

const WorkflowRepositoriesPath = "test-fixtures/user.tf"

func TestAccWorkflowRepositories_basic(t *testing.T) {
	tmpl := testAccConfigTemplate(WorkflowRepositoriesPath, t)
	name := "test"
	resourceName := "data.galaxy_workflow_repositories." + name
	workflow, _ /*parsedWorkflow*/, err := loadWorkflow(WorkflowPath)
	if err != nil {
		t.Fatal(err)
	}
	type tmplFields struct {
		name string
		json string
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		IDRefreshName:     resourceName,
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(tmpl, t, &tmplFields{name: name, json: strings.Replace(workflow, "\"", "\\\"", -1)}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "json", workflow),
					//testCheckResourceAttrLen(resourceName, "repositories", len(parsedWorkflow["steps"].(map[string]interface{}))),
				),
			},
		},
	})
}
