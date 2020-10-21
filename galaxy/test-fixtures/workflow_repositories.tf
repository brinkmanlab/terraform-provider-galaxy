data "galaxy_workflow_repositories" "{{ .Name }}" {
  json = <<EOF
{{ .Json }}
EOF
}