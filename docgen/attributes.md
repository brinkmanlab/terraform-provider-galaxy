{{ $l := .Level }}{{ range $name, $schema := .Schema }}{{ with $schema -}}
{{i $l}}* `{{ $name }}` - &lt;{{ typeName .Type }}&gt; {{ if .Deprecated }}*Depreciated* {{ end }}{{ html .Description }}  
{{ with .Elem }}{{ if isResource . -}}
{{i $l}}  Attributes:  
{{ with .Description -}}
{{i $l}}  {{ . }}  
{{ end -}}
{{ template "attributes.md" tmplParams "" .Schema .DeprecationMessage .Description (inc $l) }}
{{ with .DeprecationMessage }}{{i $l}}  !> {{ html . }}  
{{ end }}
{{- else -}}
{{i $l}}  Element type: {{ typeName .Type }}
{{ end }}{{ end -}}
{{ end }}{{ end -}}