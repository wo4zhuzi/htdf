package mint

import (
	"crypto/sha256"
	"math"
	"strconv"
)

const (
	BlkTime             = 5
	AvgDaysPerMonth     = 30
	DayinSecond         = 24 * 3600
	AvgBlksPerMonth     = AvgDaysPerMonth * DayinSecond / BlkTime
	Period              = AvgBlksPerMonth
	DefaultMineableBlks = 40 * 12 * AvgBlksPerMonth // 40 years mining
	TestMineableBlks    = 2 * 3600 / BlkTime        // 2 hours mining
	TotalMineableBlks   = DefaultMineableBlks
	BlkRadianIntv       = 2.0 * math.Pi / float64(Period)

	ValidatorNumbers         = 7                                      // the number of validators
	ValidatorProvisions      = float64(100)                           // 100 for each validator
	ValidatorTotalProvisions = ValidatorProvisions * ValidatorNumbers // 100 for each validator

	IssuerAmount = float64(1000000) // this is for test. 0 for production, 1000000 for test

	FixedMineProvision  = float64(36000000)
	MineTotalProvisions = FixedMineProvision - ValidatorTotalProvisions - IssuerAmount // ~36,000,000 for 40 years
	AnnualProvisions    = MineTotalProvisions / 40                                     // ~900000 per year
	MonthProvisions     = AnnualProvisions / 12                                        // ~75000 per month

	UserProvisions             = float64(60000000)
	CurrentProvisions          = UserProvisions + ValidatorTotalProvisions + IssuerAmount // ~60,000,000 at genesis
	CurrentProvisionsAsSatoshi = int64(CurrentProvisions * htdf2satoshi)                  // ~60,000,000 at genesis
	TotalLiquid                = MineTotalProvisions + CurrentProvisions                  // 96,000,000
	TotalLiquidAsSatoshi       = int64(TotalLiquid * htdf2satoshi)                        // 96,000,000 * 100,000,000

	htdf2satoshi = 100000000 // 1 htdf = 10 ** 8 satoshi

	// junying-todo, 2019-12-05
	AvgBlkReward          = MineTotalProvisions / TotalMineableBlks
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

func randomUint(seed int64) uint64 {
	hash := sha256.Sum256([]byte(strconv.FormatInt(seed, 10)))
	return uint64(hash[:1][0])
}

func randomFloat(seed int64) float64 {
	return float64(randomUint(seed)) / 256.0
}

// junying-todo, 2019-12-06
// rand(0.001,0.144676)
// 0.144676 * rand(0.0,1,0) + 0.001
func randomAmplitude(seed int64) int64 {
	ampf := float64(htdf2satoshi) * ((MAX_AMPLITUDE-MIN_AMPLITUDE)*randomFloat(seed) + MIN_AMPLITUDE)
	return int64(ampf)
}

// rand(100,3000)
// rand(0,2900) + 100
func randomCycle(seed int64) int64 {
	return int64(randomFloat(seed)*float64(MAX_CYCLE-MIN_CYCLE)) + MIN_CYCLE
}

//
func calcRewardFloat(amp int64, cycle int64, step int64) float64 {
	if cycle == 0 {
		return 0.0
	}
	radian := 2.0 * math.Pi * float64(step) / float64(cycle)
	return float64(amp)*math.Sin(radian) + AvgBlkRewardAsSatoshi
}

func calcRewardAsSatoshi(amp int64, cycle int64, step int64) int64 {
	return int64(calcRewardFloat(amp, cycle, step))
}
