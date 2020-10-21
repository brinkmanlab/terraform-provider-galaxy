resource "galaxy_repository" "{{ .Name }}" {
  tool_shed = "{{ .Toolshed }}"
  owner = "{{ .Owner }}"
  name = "{{ .RepoName }}"
  changeset_revision = "{{ .ChangesetRevision }}"
  remove_from_disk = true
}