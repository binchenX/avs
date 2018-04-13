{{ $spec := .}}

{{with .BoardConfig}}
BOARD_FLASH_BLOCK_SIZE := {{.PartitionTable.FlashBockSize}}
{{range .PartitionTable.Partitions }}
BOARD_{{.Name | ToUpper}}IMAGE_PARTITION_SIZE := {{.Size}}
BOARD_{{.Name | ToUpper}}IMAGE_FILE_SYSTEM_TYPE := {{.Type}}
{{end }}

{{- if . | UserImageExt4}}
TARGET_USERIMAGES_USE_EXT4 := true
{{- end}}

{{- if .Bootloader}}
TARGET_NO_BOOTLOADER := false
{{- end}}

{{- if .Bootloader.Has2ndBootloader}}
TARGET_BOOTLOADER_IS_2ND = true
{{- end}}

{{- if $spec.BootImage.Kernel}}
TARGET_NO_KERNEL := false
{{- end}}

{{- with .Target}}
{{- if .NoRecovery }}
TARGET_NO_RECOVERY := true
{{- end}}

{{- if .NoRadio }}
TARGET_NO_RADIOIMAGE := true
{{- end}}

TARGET_ARCH := {{ (index .Archs 0).Name }}
TARGET_ARCH_VARIANT := {{ (index .Archs 0).Variant }}
TARGET_CPU_VARIANT := {{ (index .Archs 0).CPU.Variant }}
TARGET_CPU_ABI := {{ (index .Archs 0).CPU.Abi }}
TARGET_CPU_ABI2 := {{ (index .Archs 0).CPU.Abi2 }}

{{if eq (len .Archs)  2}}
TARGET_2ND_ARCH := {{ (index .Archs 1).Name }}
TARGET_2ND_ARCH_VARIANT := {{ (index .Archs 1).Variant }}
TARGET_2ND_CPU_VARIANT := {{ (index .Archs 1).CPU.Variant }}
TARGET_2ND_CPU_ABI := {{ (index .Archs 1).CPU.Abi }}
TARGET_2ND_CPU_ABI2 := {{ (index .Archs 1).CPU.Abi2 }}
{{end}}

{{- if eq .Binder "64"}}
TARGET_USES_64_BIT_BINDER := true
{{- end}}

{{- if .BoardPlatform}}
TARGET_BOARD_PLATFORM := {{- .BoardPlatform -}}
{{- else}}
TARGET_BOARD_PLATFORM := {{ $spec.Product.Name}}
{{- end}}

{{- if $spec.BoardConfig.Bootloader.BoardName}}
TARGET_BOOTLOADER_BOARD_NAME := {{- $spec.BoardConfig.Bootloader.BoardName -}}
{{- else}}
TARGET_BOOTLOADER_BOARD_NAME := {{ $spec.Product.Name}}
{{- end}}

{{- end}} {{/**Target**/}}

BOARD_KERNEL_CMDLINE := {{ $spec | FullCmdLine}}
{{end }}

#sepolicy
BOARD_SEPOLICY_DIRS := {{ .BoardConfig.SELinux.PolicyDir }}

# HAL's build config
{{- range .Hals }}
{{- if .BuildConfigs}}
# build config of feature {{.Name}}
{{- range .BuildConfigs}}
{{.}}
{{- end }}
{{- end}}
{{- end}}
