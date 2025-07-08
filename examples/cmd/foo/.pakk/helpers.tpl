{{- define "binary" -}}
{{ .Project.Out }}/bin/{{ .Values.name }}
{{- end }}
