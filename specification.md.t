# Validation Specification for {{ .Project  }}


{{- $mlen := len .Markdowns }}

{{ if gt 0 $mlen }}
MEOW
{{end -}}


## Stories
{{ range .Stories }}
### {{ .Name }}
    {{$story := . }}
    {{- if .Tags }}
#### Tags:
        {{ range .Tags }}
 * {{ . }}
        {{end}}
    {{end}}

    {{ if .Tests }}
#### Tests
Name | Risk | Passed | Date
-----|------|--------|-----
{{ range .Tests -}}
{{- .Name }} | {{ $story.Risk }} | {{ .Passed }} | {{ .Date }}
        {{end}}

    {{end}}

{{end}}