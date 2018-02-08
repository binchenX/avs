{{with .BootImage.Rootfs.Fstab}}

{{- range .Mounts }}
{{.Src}}    {{.Dst}}    {{.Type}}    {{.MntFlag}}    {{.FsMgrFlag}}
{{- end}}

{{end}}