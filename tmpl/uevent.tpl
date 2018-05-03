{{with .BootImage.Rootfs.UeventRc}}

{{- range .Rules }}
{{.Node}}    {{.Attr}}    {{.Mode}}    {{.UID}}    {{.GUID}}
{{- end}}

{{end}}


{{- range .Hals}}
{{- if .Devices}}
# device nodes for HAl {{.Name}}
{{- range .Devices }}
{{.Node}}    {{.Attr}}    {{.Mode}}    {{.UID}}    {{.GUID}}
{{- end}}
{{- end -}}

{{end}}