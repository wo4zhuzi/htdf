package subspace

import (
	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
)

// Used for associating paramsubspace key and field of param structs
type ParamSetPair struct {
	Key   []byte
	Value interface{}
}

// Slice of KeyFieldPair
type ParamSetPairs []ParamSetPair

// Interface for structs containing parameters for a module
type ParamSet interface {
	ParamSetPairs() ParamSetPairs
	Validate(key string, value string) (interface{}, sdk.Error)
	GetParamSpace() string
	StringFromBytes(*codec.Codec, string, []byte) (string, error)
	String() string
}
