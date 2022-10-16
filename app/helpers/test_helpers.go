package helpers

// SimAppChainID hardcoded chainID for simulation
const (
	SimAppChainID = "simapp-1"
)

type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} {
	return nil
}
