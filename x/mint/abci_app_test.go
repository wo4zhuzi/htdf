package mint

import (
	"fmt"
	"github.com/magiconair/properties/assert"
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
	fmt.Printf("a=%s|b=%s|c=%s", a.String(), b.String(), c.String())
	assert.Equal(t, a.String() == "900000.000000000000000000", true)
	assert.Equal(t, b.String() == "900.000000000000000000", true)
	assert.Equal(t, c.String() == "14467592.000000000000000000", true)
	//fmt.Printf("a=%s|b=%s|c=%s", a.String(), b.String(), c.String())
}
