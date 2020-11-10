resource "galaxy_repository" "awkscript" {
  tool_shed = "toolshed.g2.bx.psu.edu"
  owner = "brinkmanlab"
  name = "awkscript"
  changeset_revision = "7966a43dbc9e"
  remove_from_disk = true
}

data "galaxy_tool" "{{ .Name }}" {
  depends_on = [galaxy_repository.awkscript]
  id = "{{ .Id }}"
}