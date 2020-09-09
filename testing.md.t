# Validation Testing for {{ .Project  }} {{ .Release  }}

{{ .Scope }}

## Test Candidate
Commit hashes identifying the test candidate for the relevant repositories:


{{- range .Commits }}

**{{ .Repo }}** {{ .Commit }}


{{- end }}

## Tests

Test Name | Pass | Date
----------|------|------
{{- range .Stories }}
{{- range .Tests }}
{{ .Test }} | {{ .Passed }} | {{ .Time }}
{{- end }}
{{- end }}