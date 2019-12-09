package mint

import (
	sdk "github.com/orientwalt/htdf/types"
)

// expected staking keeper
type StakingKeeper interface {
	TotalTokens(ctx sdk.Context) sdk.Int
	BondedRatio(ctx sdk.Context) sdk.Dec
	InflateSupply(ctx sdk.Context, newTokens sdk.Int)
	// junying-todo, 2019-12-06
	Amplitude(ctx sdk.Context) int64
	Cycle(ctx sdk.Context) int64
	LastIndex(ctx sdk.Context) int64
	SetAmplitude(ctx sdk.Context, amp int64)
	SetCycle(ctx sdk.Context, cycle int64)
	SetLastIndex(ctx sdk.Context, index int64)
}

// expected fee collection keeper interface
type FeeCollectionKeeper interface {
	AddCollectedFees(sdk.Context, sdk.Coins) sdk.Coins
}
