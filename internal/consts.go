package internal

const (
	downloadsUrl           string = "https://developers.google.com/android"
	StableFactoryImagesURL        = downloadsUrl + "/images"
	StableOTAImagesURL            = downloadsUrl + "/ota"
	Android17DPURL                = "https://developer.android.com/about/versions/15/download"
)

const (
	acksCookie        string = "devsite_wall_acks="
	OTAAcksCookie            = acksCookie + "nexus-ota-tos"
	FactoryAcksCookie        = acksCookie + "nexus-image-tos"
	Android17DPCookie        = acksCookie + "android-preview-tos"
)
