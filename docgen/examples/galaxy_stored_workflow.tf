resource "galaxy_stored_workflow" "example" {
  json = file("workflow.ga")
}