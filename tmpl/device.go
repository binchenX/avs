package tmpl

// Device is the template for device.mk
const Device = `# 1. bootimage
# 1.1 kernel and dtb
LOCAL_KERNEL := {{ .BootImage.Kernel.LocalKernel}}
PRODUCT_COPY_FILES += $(LOCAL_KERNEL):kernel
LOCAL_DTB := {{.BootImage.Kernel.LocalDTB}}
# TODO: fix the dest dtb name, normal it is some varient of product.dtb
PRODUCT_COPY_FILES += $(LOCAL_KERNEL):dtb
{{- if .BoardConfig.Bootloader.Has2ndBootloader }}
PRODUCT_COPY_FILES += $(LOCAL_DTB):2ndbootloader
{{- end}}

# 1.2 rootfs
PRODUCT_COPY_FILES += \
    {{ .BootImage.Rootfs.UeventRc | getUeventdCopySrc}}:root/ueventd.{{- .Product.Name -}}.rc \
    {{ .BootImage.Rootfs.Fstab    | getFstabCopySrc}}:root/fstab.{{- .Product.Name}}

PRODUCT_COPY_FILES += \
    {{ .BootImage.Rootfs.InitRc | getInitRcCopyStatement}}
{{- if .BoardConfig.USBGadget}}
    $(LOCAL_PATH)/rootfs/init.{{.Product.Name}}.usb.rc:root/init.{{.Product.Name}}.usb.rc \
{{- end}}

{{if .BoardConfig.BoardFeatures}}
# feature declaration
PRODUCT_COPY_FILES += \ 
{{- range .BoardConfig.BoardFeatures }}
    {{FeatureFileSrcDir}}/{{.}}:{{FeatureFileDestDir}}/{{.}} \
{{- end}}
{{end}}

{{- if .FrameworkConfigs}}
{{- if .FrameworkConfigs.Properties}}
# framework properties
PRODUCT_PROPERTY_OVERRIDES += \
    {{- range .FrameworkConfigs.Properties}}
    {{.}} \
    {{- end }}
{{- end}}
{{- end }}

{{range .Hals }}

# start HAL {{.Name}} >>>>>>>>
{{- if .Features}}
## feature declaration
PRODUCT_COPY_FILES += \ 
{{- range .Features }}
    {{FeatureFileSrcDir}}/{{.}}:{{FeatureFileDestDir}}/{{.}} \
{{- end}}
{{- end}}

{{with .Packages}}

{{- if .Build}}
## build packages
PRODUCT_PACKAGES += \
{{- range .Build}}
    {{. | removeTag }} \
{{- end}}
{{- end}}

{{- if .Copy}}
## copy packages
PRODUCT_COPY_FILES += \
{{- range .Copy }}
    {{ . | CopyInstruction }} \
{{- end}}
{{- end}}

{{- end -}} {{/**Packages**/}}

{{- if .Firmwares}}
## firmwares
PRODUCT_COPY_FILES += \
{{- range .Firmwares }}
    {{ . | InstsallFirmware }} \
{{- end}}
{{- end}}

{{- if .Drivers}}
## drivers
PRODUCT_COPY_FILES += \
{{- range .Drivers }}
    {{ . | InstsallDriver }} \
{{- end}}
{{- end}}

{{- if .RawInstructions}}
# raw instructions - do I have a better place to go?
{{- range .RawInstructions}}
{{.}}
{{- end}}
{{- end}}

{{- if .InitRc}}
## service init.rc scripts
PRODUCT_COPY_FILES += \
    {{ .InitRc | getInitRcCopyStatement}}
{{- end}}

{{- if .RuntimeConfigs}}
## runtime configs
PRODUCT_COPY_FILES += \
{{- range .RuntimeConfigs }}
    {{ . | RuntimeConfigInstructions}} \
{{- end}}
{{- end}}

{{- if .Properties}}
## feature {{.Name}} properties
PRODUCT_PROPERTY_OVERRIDES += \
    {{- range .Properties}}
    {{.}} \
    {{- end}}
{{- end}}

{{end}}{{/**Hals**/}}

# manifest.xml
PRODUCT_COPY_FILES += \
    $(LOCAL_PATH)/{{"manifest.xml"}}:${TARGET_COPY_OUT_VENDOR}/manifest.xml

{{if .VendorRaw}}
# vendor raw instructions - does it has a better place to go?
{{- range .VendorRaw.Instructions}}
{{.}}
{{- end}}
{{end}}`
