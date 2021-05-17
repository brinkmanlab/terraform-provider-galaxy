resource "galaxy_quota" "{{ .Name }}" {
  name = "{{ .Quotaname }}"
  description = "test"
  default = "{{ .DefaultFor }}"
  amount = "{{ .Amount }}"
}