{{ $l := .Level -}}
{{ range $name, $schema := .Schema }}{{ with $schema }}{{ if or .Optional (not .Computed) -}}
{{i $l}}* `{{ $name }}` - &lt;{{ typeName .Type }}&gt; {{ if .Deprecated }}*Depreciated* {{ end }}{{ if .Required }}(Required) {{ end }}{{ if .Optional }}(Optional) {{ end }}{{ html .Description }}{{ with .Default }} \[Default: {{ html . }}]{{ end }}  
{{ with .ExactlyOneOf -}}
{{i $l}}  Exactly one of {{ listAttr . "or" }}  
{{ end -}}
{{ with .ConflictsWith -}}
{{i $l}}  Conflicts with {{ listAttr . "and" }}  
{{ end -}}
{{ with .AtLeastOneOf -}}
{{i $l}}  At least one of {{ listAttr . "or" }}  
{{ end -}}
{{ with .RequiredWith -}}
{{i $l}}  Required with {{ listAttr . "and" }}  
{{ end -}}
{{- if or (gt .MinItems 0) (gt .MaxItems 0) }}
{{i $l}}  Limit {{ .MinItems }}-{{ if gt .MaxItems 0 }}{{ .MaxItems }}{{ else }}&#8734;{{ end }} items  
{{ end -}}
{{- with .Elem }}{{ if isResource . -}}
{{i $l}}  Arguments:  
{{ with .Description -}}
{{i $l}}  {{ html . }}  
{{ end -}}
{{ template "arguments.md" tmplParams "" .Schema .DeprecationMessage .Description (inc $l) }}
{{ with .DeprecationMessage }}{{i $l}}  !> {{ html . }}  
{{ end }}
{{- else -}}
{{i $l}}  Element type: {{ typeName .Type }}
{{ end }}{{ end -}}
{{ end }}{{ end }}{{ end -}}