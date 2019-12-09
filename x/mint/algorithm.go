package mint

import (
	"math"
	"math/rand"
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

	htdf2satoshi = 100000000 // 1 htdf = 10 ** 8 satoshi

	// junying-todo, 2019-12-05
	AvgBlkReward          = MonthProvisions / AvgBlksPerMonth
	AvgBlkRewardAsSatoshi = htdf2satoshi * AvgBlkReward
	RATIO                 = 0.5
	// junying-todo, 2019-12-06
	MAX_AMPLITUDE = AvgBlkReward
	MIN_AMPLITUDE = 0.001
	MAX_CYCLE     = 3000
	MIN_CYCLE     = 100
)

// junying-todo, 2019-12-05
// 60,000,000 + 0.144676 * height
func expectedtotalSupply(blkindex int64) int64 {
	return CurrentProvisionsAsSatoshi +
		int64(float64(blkindex)*AvgBlkRewardAsSatoshi)
}

// junying-todo, 2019-12-06
// rand(0.001,0.144676)
// 0.144676 * rand(0.0,1,0) + 0.001
func randomAmplitude() int64 {
	ampf := float64(htdf2satoshi) * (MAX_AMPLITUDE*rand.Float64() + MIN_AMPLITUDE)
	return int64(ampf)
}

// rand(100,3000)
// rand(0,2900) + 100
func randomCycle() int64 {
	return rand.Int63n(MAX_CYCLE-MIN_CYCLE) + MIN_CYCLE
}

//
func calcReward(amp int64, cycle int64, step int64) float64 {
	radian := 2.0 * math.Pi * float64(step) / float64(cycle)
	return float64(amp)*math.Sin(radian) + AvgBlkReward
}

func calcRewardAsSatoshi(amp int64, cycle int64, step int64) int64 {
	return int64(calcReward(amp, cycle, step))
}
