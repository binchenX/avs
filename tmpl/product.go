package tmpl

// Product is the template for $(Product).mk
const Product = `{{with .Product}}

{{- if .InheritProducts}}
{{- range .InheritProducts}}
$(call inherit-product, {{ . | InheritProduct }})
{{- end}}
{{- end}}

PRODUCT_NAME := {{.Name}}
PRODUCT_DEVICE := {{.Device}}
PRODUCT_BRAND := {{.Brand}}
PRODUCT_MODEL := {{.Model}}
PRODUCT_MANUFACTURER := {{.Manufacture}}

DEVICE_PACKAGE_OVERLAYS := device/{{- .Manufacture -}}/{{- .Device -}}/overlay

# automatically called
$(call inherit-product, device/{{- .Manufacture -}}/{{- .Device -}}/device.mk)

{{end}}`
