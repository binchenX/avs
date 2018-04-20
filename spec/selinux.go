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

// To audit the selinux warnings
// see https://source.android.com/security/selinux/validate
// not use audit2allow comes with your distro but the one in
// external/selinux/prebuilts/bin/audit2allow,
// Otherwise, you will mean error:
// libsepol.policydb_read: policydb version 30 does not match my version range 15-29
// invalid binary policy policy.
// Also, running the command in ANDRIOD_BUILD_TOP otherwise there will be link issue.
// You can either pull the selinux policy from the board
// adb pull /sys/fs/selinux/policy
// or use the one in the out/target/product/{devicename}/root/sepolicy,
// providing that is the same actually being used in the device
