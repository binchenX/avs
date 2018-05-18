package tmpl

// Androidproducts is the template for vendorsetup.mk
const Androidproducts = `
{{- with .Product}}
PRODUCT_MAKEFILES := \
	$(LOCAL_DIR)/{{- .Name -}}.mk
{{- end}}
`
