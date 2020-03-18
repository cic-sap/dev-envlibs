{{- define "name" -}}
{{- .Chart.Name -}}
{{- end -}}

{{- define "fullname" -}}
{{- .Chart.Name -}}
{{- end -}}

{{- define "domain" -}}
{{- printf .Values.ingress.host .Release.Namespace .Values.domain -}}
{{- end -}}
