{{- range .Imports }}
import {{.}}
{{- end}}

{{range .Actions }}
on {{.Triggers}}
    {{- range .Commands}}
        {{.}}
    {{- end}}
{{end}}


{{range .Services}}
service {{.Name}} {{.Path}} {{.Args}}
   {{- range .Options}}
    {{.}}
   {{- end }}
{{end }}
