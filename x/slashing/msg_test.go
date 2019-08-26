package slashing

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/orientwalt/htdf/types"
)

func TestMsgUnjailGetSignBytes(t *testing.T) {
	addr := sdk.AccAddress("abcd")
	msg := NewMsgUnjail(sdk.ValAddress(addr))
	bytes := msg.GetSignBytes()
	require.Equal(t, string(bytes), `{"address":"cosmosvaloper1v93xxeqhg9nn6"}`)
}
