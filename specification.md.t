# Validation Specification for {{ .Project  }}


{{ if .Markdown }}
    {{- range .Markdown }}
        {{- .Content }}
    {{- end -}}
{{- end }}



{{ range .Releases }}
## Release {{ .Name }}
### {{ .Scope }}

{{ if .Markdown }}
    {{- range .Markdown}}
        {{- .Content }}
    {{end -}}
{{end -}}


{{ end -}}