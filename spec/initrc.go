package spec

// RcScripts is a script follow Android init syntax and semantics, see[1].
// [1]https://android.googlesource.com/platform/system/core/+/master/init/README.md

// RcScripts support two modes, `External` and `Embed-in`.
// In External mode, the rc scripts isn't managed by the avs but only referred to,
// using `File` attribute. The `File` is the file path relative to the device configration
// directory. In external mode, all the other attributions will be ignored.
// In Embed-In mode, all rc scripts statement will be in the avs config file. Avs will automatically
// generate a rc scripts in the path specified by `Name` attribute. The `Name` is the file path
// relative to the device configration directory.
// The Embed-In mode is encouraged to be used for ServiceRc, i.e when SercieRc is true.
// While the external mode can be used when for init.rc in the rootfs which contains more stuff
// that deserve an external file.
type RcScripts struct {
	// ServicRc indicates if this rc is for service only and it is usually used for
	// start a service for a specific hal. Default to false.
	// The differnce is the install destination:
	// serviceRc:  "$(TARGET_COPY_OUT_VENDOR)/etc/init/",
	// initRc:    "root/"
	ServicRc string      `json:"serviceRc,omitempty"`
	File     string      `json:"file,omitempty"`
	Name     string      `json:"name,omitempty"`
	Imports  []string    `json:"imports,omitempty"`
	Actions  []RcAction  `json:"actions,omitempty"`
	Services []RcService `json:"services,omitempty"`
}

// RcAction is the Action statement.
type RcAction struct {
	Triggers string   `json:"triggers"`
	Commands []string `json:"commands"`
}

// RcService is the Service statement.
type RcService struct {
	Name    string   `json:"name"`
	Path    string   `json:"path"`
	Args    string   `json:"args,omitempty"`
	Options []string `json:"options,omitempty"`
}

// RcImport is the Import statement.
type RcImport struct {
	ImportPath string `json:"path"`
}
