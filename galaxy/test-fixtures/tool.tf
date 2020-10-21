resource "galaxy_repository" "awkscript" {
  tool_shed = "toolshed.g2.bx.psu.edu"
  owner = "brinkmanlab"
  name = "awkscript"
  changeset_revision = "ceac6ffb3865"
  remove_from_disk = true
}

data "galaxy_tool" "example" {
  depends_on = [galaxy_repository.awkscript]
  id = "toolshed.g2.bx.psu.edu/repos/brinkmanlab/awkscript/awkscript/1.0"  # An issue with the Galaxy API requires this. https://github.com/galaxyproject/galaxy/issues/10378
}