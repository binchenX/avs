package tmpl

// Vendorsetup is the template for vendorsetup.mk
const Vendorsetup = `
{{with .Product}}
add_lunch_combo {{ .Name -}}-eng
{{end}}`
