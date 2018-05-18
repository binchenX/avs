package tmpl

// Fstab is the template for fstab.$(product)
const Fstab = `{{with .BootImage.Rootfs.Fstab}}

{{- range .Mounts }}
{{.Src}}    {{.Dst}}    {{.Type}}    {{.MntFlag}}    {{.FsMgrFlag}}
{{- end}}

{{end}}`
