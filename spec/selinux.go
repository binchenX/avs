package spec

// SELinux setting for the board
// https://source.android.com/security/selinux/

const (
	// SElinuxModePermissive means permission denials are logged but not enforced
	SElinuxModePermissive = "permissive"
	// SElinuxModeEnforcing means permissions denials are both logged and enforced.
	SElinuxModeEnforcing = "enforcing"
)

// SELinux setting
type SELinux struct {
	// Themde set here will override the kernel commandline and default is SElinuxModeEnforcing
	Mode      string `json:"mode"`
	PolicyDir string `json:"policyDir"`
}

// SEPolicyF is the sepolicy configration.
type SEPolicyF struct {
	// create new file type
	FileTe []string `json:"file.te,omitempty"`
	// create new process domain
	ServiceTe []string `json:"service.te,omitempty"`
	// lable the files
	FileContexts []string `json:"file_contexts,omitempty"`
	// lable the proccess/services
	ServiceContexts []string `json:"service_contexts,omitempty"`
}
