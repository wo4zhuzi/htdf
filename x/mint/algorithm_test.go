package mint

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestCalcMiningReward(t *testing.T) {
	var sum float64
	for i := 0; i < TotalMinableBlks; i++ {
		sum += calcMiningReward(i)
	}
	fmt.Println("sum:", sum)
	require.Equal(t, int(math.Round(sum)), TotolProvisions)
}

func TestEstimatePeak(t *testing.T) {
	estimated := estimatePeak()
	fmt.Println("RewardsPerMonth:", MonthProvisions)
	fmt.Println("AvgBlksPerMonth:", AvgBlksPerMonth)
	fmt.Println("estimatedPeak:", estimated)
	require.Equal(t, estimated, Peak)
}
