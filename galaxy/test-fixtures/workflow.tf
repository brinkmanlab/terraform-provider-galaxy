resource "galaxy_stored_workflow" "{{ .Name }}" {
  json = <<EOF
{{ .Json -}}
EOF
}