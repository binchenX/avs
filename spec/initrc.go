package spec

// RcScripts is a script follow Android init syntax and semantics, see[1].
// [1]https://android.googlesource.com/platform/system/core/+/master/init/README.md
type RcScripts struct {
	// If Files is no nil, we will cp it directly $(LOCAL_PATH)/File to destination
	// All the other attribution will be ignored
	File string `json:"file,omitempty"`
	// Name will be the file name for the generated script and it will be copied to device
	// with same name. At least one of the File, or Name should be non empty.
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
