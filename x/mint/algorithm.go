package mint

import (
	"math"
)

const (
	BlkTime          = 5
	AvgDaysPerMonth  = 30
	DayinSecond      = 24 * 3600
	AvgBlksPerMonth  = AvgDaysPerMonth * DayinSecond / BlkTime
	Period           = AvgBlksPerMonth
	TotalMinableBlks = 40 * 12 * AvgBlksPerMonth
	BlkRadianIntv    = 2.0 * math.Pi / float64(Period)

	MonthProvisions            = float64(75000)                          // 75000 per month
	AnnualProvisions           = 12 * MonthProvisions                    // 900000 per year
	TotalProvisions            = 40 * AnnualProvisions                   // 36,000,000 for 40 years
	CurrentProvisions          = float64(60000000)                       // 60,000,000 at genesis
	CurrentProvisionsAsSatoshi = int64(CurrentProvisions * htdf2satoshi) // 60,000,000 at genesis
	TotalLiquid                = TotalProvisions + CurrentProvisions     // 96,000,000
	TotalLiquidAsSatoshi       = int64(TotalLiquid * htdf2satoshi)       // 96,000,000 * 100,000,000

	Peak         = 0.1446759259259267
	DeltaRadian  = 0.0
	htdf2satoshi = 100000000 // 1 htdf = 10 ** 8 satoshi

	// junying-todo, 2019-12-05
	AvgBlkReward          = MonthProvisions / AvgBlksPerMonth
	AvgBlkRewardAsSatoshi = htdf2satoshi * AvgBlkReward

	RATIO = 0.5
)

// junying-todo, 2019-12-05
func estimatedAccumulatedSupply(blkindex int64) int64 {
	return CurrentProvisionsAsSatoshi +
		int64(float64(blkindex)*AvgBlkRewardAsSatoshi)
}
