<manifest version="1.0" type="device">
{{- range .Hals -}}
{{- if .Manifests}}
{{- range .Manifests}}
    <hal format="{{- .Format -}}">
        <name>{{- .Name -}}</name>
        {{- if .Transport.Arch }}
        <transport arch="{{.Transport.Arch}}">{{- .Transport.Mode -}}</transport>
        {{- else }}
        <transport>{{- .Transport.Mode -}}</transport>
        {{- end }}
        {{- if .Impl }}
        <impl level="{{.Impl.Level}}"></impl>
        {{- end}}
        <version>{{- .Version -}}</version>
        <interface>
            <name>{{ .Interface.Name }}</name>
            <instance>{{ .Interface.Instance }}</instance>
        </interface>
    </hal>
{{- end}}
{{- end}}
{{- end}}
</manifest>
