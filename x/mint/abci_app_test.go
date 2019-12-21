package mint

import (
	"fmt"
	"testing"

	sdk "github.com/orientwalt/htdf/types"
)

func TestMineToken(t *testing.T) {

	var curBlkHeight int64
	var curAmplitude int64
	var curCycle int64
	var curLastIndex int64
	totalSupply := sdk.NewInt(1000)

	curBlkHeight = 1
	curAmplitude = 2
	curCycle = 3
	curLastIndex = 4

	a, b, c := GetMineToken(curBlkHeight, totalSupply, curAmplitude, curCycle, curLastIndex)
	fmt.Printf("a=%v|b=%v|c=%v", a, b, c)
}
