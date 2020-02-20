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

	annualProvisions, inflation, blockProvision := GetMineToken(curBlkHeight, totalSupply, curAmplitude, curCycle, curLastIndex)
	fmt.Printf("a=%s|b=%s|c=%s\n", annualProvisions.String(), inflation.String(), blockProvision.String())
	fmt.Println(annualProvisions.ToUint64())
	assert.Equal(t, annualProvisions.ToUint64() == uint64(AnnualProvisionAsSatoshi), true)
	assert.Equal(t, inflation.String() == "90000000000.000000000000000000", true)
	assert.Equal(t, blockProvision.String() == "14467592.000000000000000000", true)
	//fmt.Printf("a=%s|b=%s|c=%s", a.String(), b.String(), c.String())
}
