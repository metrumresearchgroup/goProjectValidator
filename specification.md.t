# Validation Specification for {{ .Project  }}


{{ if .MarkDown }}
    {{- range .MarkDown }}
        {{- .Content }}
    {{- end -}}
{{- end }}


## Stories
{{ range .Stories }}

{{ if .MarkDown }}
    {{- range .MarkDown}}
        {{- .Content }}
    {{end -}}
{{end -}}


{{ end -}}