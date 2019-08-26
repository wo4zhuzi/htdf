package accounts

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/orientwalt/htdf/types"
)

const (
	// AddrLen defines a valid address length
	AddrLen = 40
)

// AccAddressFromStr creates an AccAddress from a xxx string.
func AccAddressFromStr(address string) (addr sdk.AccAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return sdk.AccAddress{}, nil
	}

	fmt.Print(len(strings.TrimSpace(address)))
	if len(strings.TrimSpace(address)) != AddrLen {
		return nil, errors.New("Incorrect address length")
	}

	return []byte(address), nil
}
