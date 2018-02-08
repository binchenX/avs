{{- with .Product}}
PRODUCT_MAKEFILES := \
	$(LOCAL_DIR)/{{- .Name -}}.mk
{{- end}}
