package mint

import (
	"fmt"
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
			curAmplitude = randomAmplitude(curBlkHeight)
			curCycle = randomCycle(curAmplitude)
			curLastIndex = curBlkHeight
		}
		BlockReward := calcRewardAsSatoshi(curAmplitude, curCycle, curBlkHeight-curLastIndex)
		// junying-todo, 2019-12-20
		// avoid negative rewards
		if BlockReward < 0 {
			break
		}
		totalSupply += BlockReward
	}
	return math.Abs(float64(curBlkHeight - lastblkindex))
}

func TestRandomSine(t *testing.T) {
	threshold := float64(MAX_CYCLE / 2)
	actual := calcErrRate(TotalMineableBlks)
	require.True(t, actual < threshold)
}

func TestRandomUint(t *testing.T) {
	cnt := 0
	for i := 0; i < 256; i++ {
		randnum := randomUint(int64(i))
		if randnum == 0 {
			cnt++
		}
		require.True(t, randnum < 256)
		require.True(t, randnum == randomUint(int64(i)))
	}
	fmt.Println(cnt)
}
