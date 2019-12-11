package mint

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func calcErrRate(lastblkindex int64) float64 {
	var totalSupply = CurrentProvisionsAsSatoshi
	var curBlkHeight int64
	var curAmplitude, curCycle, curLastIndex int64
	for curBlkHeight = 1; curBlkHeight < (lastblkindex + MAX_CYCLE/2); curBlkHeight++ {
		// check if mined is greater than expected
		if totalSupply >= TotalLiquidAsSatoshi {
			break
		}
		// check if it's time for new cycle
		if curBlkHeight >= (curLastIndex + curCycle) {
			curAmplitude = randomAmplitude()
			curCycle = randomCycle()
			curLastIndex = curBlkHeight
		}
		BlockReward := calcRewardAsSatoshi(curAmplitude, curCycle, curBlkHeight)
		totalSupply += BlockReward
	}
	return math.Abs(float64(curBlkHeight - lastblkindex))
}

func TestRandomSine(t *testing.T) {
	threshold := float64(MAX_CYCLE / 2)
	actual := calcErrRate(TotalMinableBlks)
	require.True(t, actual < threshold)
}
