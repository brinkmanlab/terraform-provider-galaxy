resource "galaxy_user" "{{ .Name }}" {
  username = "{{ .Username }}"
  password = "{{ .Password }}"
  email = "{{ .Email }}"
  purge = false
}

resource "galaxy_quota" "{{ .Name }}" {
  name = "{{ .Quotaname }}"
  description = "test"
  amount = "{{ .Amount }}"
  users = [galaxy_user.{{ .Name }}.id]
}