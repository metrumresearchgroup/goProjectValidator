# Traceability Matrix: {{ .Project }}
{{- range .Releases }}
# Release {{ .Name }}
{{ .Scope }}

Issue Title | Risk | Test Name | Pass | Date
------------|------|-----------|------|------

{{- range .Stories }}
{{- $story := . }}
{{- range .Tests }}
{{ $story.Name }} | {{ $story.Risk}} | {{ .Test }} | {{ .Passed }} | {{ .Time }}
{{- end }}
{{- end }}
{{- end }}