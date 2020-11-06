resource "galaxy_repository" "awkscript" {
  tool_shed = "toolshed.g2.bx.psu.edu"
  owner = "brinkmanlab"
  name = "awkscript"
  changeset_revision = "ceac6ffb3865"
  remove_from_disk = true
}

resource "galaxy_history" "test" {
  name = "test"
}

resource "galaxy_job" "example" {
  depends_on = [galaxy_repository.awkscript]
  tool_id = galaxy_repository.awkscript.tools[0].tool_guid
  history_id = galaxy_history.test.id
  params = {
    "code" = "BEGIN { print \"foo\" }"
  }
  wait_for_completion = true
}