data "galaxy_workflow_repositories" "example" {
  json = file("workflow.ga")
}

resource "galaxy_repository" "example" {
  for_each = data.galaxy_workflow_repositories.example.repositories
  tool_shed = each.value.tool_shed
  owner = each.value.owner
  name = each.value.name
  changeset_revision = each.value.changeset_revision
  remove_from_disk = true
}