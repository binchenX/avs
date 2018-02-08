{{with .BootImage.Rootfs.UeventRc}}

{{- range .Rules }}
{{.Node}}    {{.Attr}}    {{.Mode}}    {{.UID}}    {{.GUID}}
{{- end}}

{{end}}