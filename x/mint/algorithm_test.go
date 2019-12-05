package mint

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func calcErrRate(lastblkindex int64) float64 {
	var actual, expected int64
	for blkindex := int64(1); blkindex < lastblkindex; blkindex++ {
		expected = estimatedAccumulatedSupply(blkindex)
		var real, estimated int64
		if expected > actual {
			estimated = expected - actual
			real = rand.Int63n(int64(float64(estimated)*RATIO) + estimated)
		}
		actual += real
	}
	return float64(actual-expected) / htdf2satoshi
}

func TestSimulate(t *testing.T) {
	threshold := float64(AvgBlkReward * 4)
	require.True(t, math.Abs(calcErrRate(1001)) < threshold)
	require.True(t, math.Abs(calcErrRate(10001)) < threshold)
	require.True(t, math.Abs(calcErrRate(100001)) < threshold)
	require.True(t, math.Abs(calcErrRate(1000001)) < threshold)
	require.True(t, math.Abs(calcErrRate(10000001)) < threshold)
	require.True(t, math.Abs(calcErrRate(TotalMinableBlks)) < threshold)
}
