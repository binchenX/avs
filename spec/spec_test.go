package spec

func ExampleSpec() {
	_ = Spec{
		Version: &Version{
			Schema:  "0.1",
			Android: "Android o",
		},

		Product: &Product{
			Name:        "poplar",
			Device:      "Poplar",
			Brand:       "Poplar",
			Model:       "Poplar",
			Manufacture: "linaro",
		},
	}
}

func ExampleProduct() {
	_ = &Product{
		Name:            "poplar",
		Device:          "Poplar",
		Brand:           "Poplar",
		Model:           "Poplar",
		Manufacture:     "linaro",
		InheritProducts: []string{"full_base"},
	}
}

// This example will label a vendor device node as type uhid_device, which has already
// being defined by Android.
func ExampleSEPolicyF_labelANewFile() {
	_ = SEPolicyF{
		FileContexts: []string{
			"/dev/stpbt	 u:object_r:uhid_device:s0"},
	}
}

// This example will show you 1) how to define a *new* label "sync_file",
// 2) use it to label a device node (which may have other labels), and
// 3) and allows a daemon (which *already* in a domain) to access it.
func ExampleSEPolicyF_createANewLabel() {
	_ = &SEPolicyF{
		FileTe: []string{
			"type sync_file, fs_type, debugfs_type;",
		},
		FileContexts: []string{
			"/sys/kernel/debug/sync  u:object_r:sync_file:s0",
		},
		ServiceTe: []string{
			"allow surfaceflinger sync_file:file rw_file_perms;",
		},
	}
}

// This example will show you what need to be done for a new daemon process.
func ExampleSEPolicyF_labeANewDaemon() {
	_ = &SEPolicyF{
		FileContexts: []string{
			"/system/vendor/bin/vendor_process   u:object_r:vendor_processr_exec:s0",
		},
		ServiceTe: []string{
			// convenstionly, it should be vendor_process.te
			"type vendor_process, domain;",
			"type vendor_processr_exec, exec_type, file_type;",
			"init_daemon_domain(vendor_process)",
			"binder_use(vendor_process)",
			"binder_service(vendor_process)",
			"binder_call(vendor_process, system_server)",
			"allow vendor_process graphics_device:dir search;",
			// convenstionly, it should be service.te
			// https://android.googlesource.com/platform/system/sepolicy/+/master/public/service.te
			// but it's ok to put them together
			"type vendor_process_service,  service_manager_type;",
		},
		ServiceContexts: []string{
			"vendor_process    u:object_r:vendor_process_service:s0",
		},
	}
}
