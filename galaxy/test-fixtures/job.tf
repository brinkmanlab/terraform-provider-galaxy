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

resource "galaxy_job" "{{ .Name }}" {
  depends_on = [galaxy_repository.awkscript]
  tool_id = "toolshed.g2.bx.psu.edu/repos/brinkmanlab/awkscript/awkscript/1.0"
  history_id = galaxy_history.test.id
  params = {
    "code" = "BEGIN { print \"foo\" }"
  }
  wait_for_completion = true
}