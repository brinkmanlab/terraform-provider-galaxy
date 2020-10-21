package galaxy_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

const WorkflowRepositoriesPath = "./test-fixtures/workflow_repositories.tf"

func TestAccWorkflowRepositories_basic(t *testing.T) {
	tmpl := testAccConfigTemplate(WorkflowRepositoriesPath, t)
	name := "test"
	resourceName := "data.galaxy_workflow_repositories." + name
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
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(tmpl, t, &tmplFields{Name: name, Json: workflow}),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckResourceAttrEqual(resourceName, "repositories.#", len(parsedWorkflow["steps"].(map[string]interface{}))),
					resource.TestCheckResourceAttr(resourceName, "repositories.0.name", "awkscript"),
				),
			},
		},
	})
}
