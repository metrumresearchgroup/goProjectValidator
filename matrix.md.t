# Traceability Matrix for {{ .Project }} {{ .Release  }}

{{ .Scope }}

Issue Title | Risk | Test Name | Pass | Date
------------|------|-----------|------|------
{{- range .Stories }}
{{- $story := . }}
{{- range .Tests }}
{{ $story.Name }} | {{ $story.Risk}} | {{ .Test }} | {{ .Passed }} | {{ .Time }}
{{- end }}
{{- end }}