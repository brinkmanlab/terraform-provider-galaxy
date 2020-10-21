# {{ .Name }} Data Source

{{ html .Description }}

{{- with .DepreciationMessage }}!> {{ html . }}{{ end }}

## Example Usage

```hcl
{{ example .Name }}
```

## Argument Reference

{{ template "arguments.md" . }}

## Attribute Reference

{{ template "attributes.md" . }}
