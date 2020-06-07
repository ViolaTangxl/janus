package config

// RunMode Run mode of app
type RunMode string

// IsProd Check if running in product mode
func (r RunMode) IsProd() bool {
	return r == ProdMode
}

var (
	// ProdMode Product mode
	ProdMode RunMode = "release"
	// DevMode Develop mode
	DevMode RunMode = "debug"
	// TestMode Test mode
	TestMode RunMode = "test"
)
