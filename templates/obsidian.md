#végétarian-{{- if .Vegetarian }}oui{{- else }}non{{- end }} #réchauffable-{{- if .Reheat }}oui{{- else }}non{{- end }} #difficulté-{{.Difficulty}} #coût-{{.Pricing}}

## Durée

- Préparation: {{ .Preparation.TimeDuration  }}
- Repos: {{ .Resting.TimeDuration  }}
- Cuisson: {{ .Cooking.TimeDuration  }}
{{ if (gt .Servings.Quantity 0) }}
## Ingrédients ({{.Servings.Quantity}} 🍽️)
{{- else }}
## Ingrédients
{{- end }}
{{ range $ingredient := .Ingredients }}
- [[{{ $ingredient.Title }}]]{{ if gt $ingredient.Quantity.Value 0.0 }} | {{ $ingredient.Quantity }}{{ end }}{{ if ne $ingredient.Details "" }} | {{ $ingredient.Details }}{{ end }}
{{- end }}

## Étapes
{{  range $instruction := .Instructions -}}
{{- if $instruction.Recipe }}
[[{{$instruction.Recipe.Title}}|{{ $instruction.Title }}]]
{{- else }}
{{ $instruction.Title }}
{{- end }}
{{- if $instruction.Warning }}=={{ $instruction.Warning }}=={{ end }}
{{ end -}}
