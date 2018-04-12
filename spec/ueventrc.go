package spec

// UeventRc is the rules for eventd.
type UeventRc struct {
	// If Files is no nil, we will cp it directly $(LOCAL_PATH)/File to destination
	// All the other attribution will be ignored
	File string `json:"file,omitempty"`
	// Name will be the file name for the generated script and it will be copied to device
	// with same name. Default name ueventd.rc.gen
	Name  string       `json:"name,omitempty"`
	Rules []UeventRule `json:"rules"`
}

// UeventRule is rule for eventd.
type UeventRule struct {
	Node string `json:"node"`
	Attr string `json:"attr,omitempty"`
	Mode string `json:"mode"`
	UID  string `json:"uid"`
	GUID string `json:"guid"`
}
