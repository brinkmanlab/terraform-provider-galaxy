resource "galaxy_user" "{{ .Name }}" {
  username = "{{ .Username }}"
  password = "{{ .Password }}"
  email = "{{ .Email }}"
}