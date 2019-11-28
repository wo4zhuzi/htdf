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
	MonthProvisions  = 75000
	AnnualProvisions = 12 * MonthProvisions
	TotolProvisions  = 40 * AnnualProvisions
	Peak             = 0.1446759259259267
	DeltaRadian      = 0.0
	Htdf2Satoshi     = 100000000
)

func estimatePeak() float64 {
	sum := 0.0
	var i int64
	for i = 0; i < Period; i++ {
		radian := 2.0 * math.Pi * float64(i) / float64(Period)
		sum += math.Sin(radian) + 1.0
	}
	return float64(75000) / sum
}

func calcMiningReward(blkheight int) float64 {
	if blkheight > TotalMinableBlks {
		return 0.0
	}
	radian := BlkRadianIntv * float64(blkheight%Period)
	return Peak * (math.Sin(DeltaRadian+radian) + 1.0)
}
