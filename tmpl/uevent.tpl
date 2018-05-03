{{with .BootImage.Rootfs.UeventRc}}
{{- range .Rules }}
{{printf "%-30s" .Node}} {{.Attr}}  {{.Mode}}  {{printf "%-10s" .UID}} {{.GUID}}
{{- end}}
{{end}}

{{- range .Hals}}
{{- if .Devices}}
# device nodes for HAl {{.Name}}
{{- range .Devices }}
{{printf "%-30s" .Node}} {{.Attr}}  {{.Mode}}  {{printf "%-10s" .UID}} {{.GUID}}
{{- end}}
{{- end -}}
{{end}}