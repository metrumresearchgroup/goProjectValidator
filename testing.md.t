# Validation Specification for {{ .Project  }}


{{ if .Markdown }}
    {{- range .Markdown }}
        {{- .Content }}
    {{- end -}}
{{- end }}


## Stories
{{ range .Releases}}
{{ range .Stories }}

{{ if .Markdown }}
    {{- range .Markdown}}
        {{- .Content }}
    {{end -}}
{{end }}

##### Test Results

Test Name | Test Output
----------|-----------
    {{- range .Tests }}
{{ .Test }} | {{ .Output }}
    {{- end}}


{{ end -}}
{{- end}}
