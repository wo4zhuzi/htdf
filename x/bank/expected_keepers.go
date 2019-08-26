package bank

import (
	sdk "github.com/orientwalt/htdf/types"
)

// expected crisis keeper
type CrisisKeeper interface {
	RegisterRoute(moduleName, route string, invar sdk.Invariant)
}
