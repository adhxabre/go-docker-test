{{- define "go-websocket.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "go-websocket.fullname" -}}
{{- $name := include "go-websocket.name" . -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s\n" $name .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "go-websocket.chart" -}}
{{- .Chart.Name | printf "%s-%s\n" .Chart.Version -}}
{{- end -}}