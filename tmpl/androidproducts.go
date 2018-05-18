package tmpl

// Androidproducts is the template for AndroidProducts.mk
const Androidproducts = `
{{- with .Product}}
PRODUCT_MAKEFILES := \
	$(LOCAL_DIR)/{{- .Name -}}.mk
{{- end}}
`
