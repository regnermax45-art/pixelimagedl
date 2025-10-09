package pixelimagedl

//go:generate stringer -type PortingAlgorithm -linecomment -output port_types_string.go

type PortingAlgorithm uint8

const (
	BasicAlgorithm       PortingAlgorithm = iota // basic
	AdvancedAlgorithm                            // advanced
	ExperimentalAlgorithm                        // experimental
)

type PortingConfig struct {
	SourceDevice     Pixel
	TargetDevice     Pixel
	Algorithm        PortingAlgorithm
	DownloadType     DownloadType
	OutputDirectory  string
}

type PortingResult struct {
	SourceImage     PixelImage
	PortedImagePath string
	Algorithm       PortingAlgorithm
	Modifications   []string
	Warnings        []string
}

type CustomFirmwareRequirements struct {
	CompatibilityLevel    string
	SelectionReason       string
	CustomModifications   []string
	ModificationCount     int
}
