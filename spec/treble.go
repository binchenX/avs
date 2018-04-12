package spec

// This file contains all things related with treble: hal manifest, vndk, dtb etc.

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
