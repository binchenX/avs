package tmpl

// Vendorsetup is the template for vendorsetup.mk
const Fstab = `{{with .BootImage.Rootfs.Fstab}}

{{- range .Mounts }}
{{.Src}}    {{.Dst}}    {{.Type}}    {{.MntFlag}}    {{.FsMgrFlag}}
{{- end}}

{{end}}`
