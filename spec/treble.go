package spec

// This file contains all things related with treble: hal manifest, vndk, dtb etc.

// https://source.android.com/devices/architecture/vintf/objects
// validate hal.format
const (
	HIDL   string = "hidl"
	NATIVE string = "native"
)

// valid hal.transport
const (
	HB string = "hwbinder"
	PT string = "passthrough"
)

// valid hal.transport.arch
const (
	A32 string = "32"
	A64 string = "64"
	AB  string = "32+64"
)

// Manifest is the manifest for the interface.
type Manifest struct {
	Name      string           `json:"name"`
	Format    string           `json:"format"`
	Transport *Transport       `json:"transport"`
	Impl      *Impl            `json:"impl,omitempty"`
	Version   string           `json:"version"`
	Interface *ServiceInterace `json:"interface"`
}

// Transport is the transport type could be hwbinder, passthrough.
type Transport struct {
	Arch string `json:"arch,,omitempty"`
	Mode string `json:"mode"`
}

// Impl is the implementation.
type Impl struct {
	Level string `json:"level"`
}

// ServiceInterace is the interface this HAL model implemented.
type ServiceInterace struct {
	Name     string `json:"name"`
	Instance string `json:"instance"`
}
