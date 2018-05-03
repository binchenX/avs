package spec

// USBGadget is the config needed for adb to work
type USBGadget struct {
	Serialnumber string `json:"serialnumber"`
	Manufacturer string `json:"manufacturer"`
	Product      string `json:"product"`
	Controller   string `json:"controller"`
}
