package mint

import (
	"fmt"

	sdk "github.com/orientwalt/htdf/types"
)

// junying-todo, 2019-07-17
//	6,000,000	25
//	6,000,000	12.5
// 	6,000,000	6.25
//	...
// ex:
//	BlksPerRound = 100
//	rewards+commission+community-pool
//	hscli query distr rewards htdf1zulqmaqlsgrgmagenaqf02p8kfgsuqkdwgwj80
//	121793749706.0satoshi
//	* 4 = 487174998824
//	not true becasue proper get more rewards,that's, different rewards on every node.
//  hscli query distr commission cosmosvaloper1lwjmdnks33xwnmfayc64ycprww49n33mtm92ne
// 	hscli query distr community-pool
const (
	// Block Reward of First Round
	InitialReward = 25 * 100000000 //25htdf = 2500000000satoshi
	// Block Count Per Round
	BlksPerRound = 6000000 //10 //6,000,000
	// Last Round Index with Block Rewards
	LastRoundIndex = 31
)

// junying-todo, 2019-07-15
// single node: 88.2 for delegators, 11.8 for validator(commission)
// commission is validating fee
// commission rate changes?
func calcParams(ctx sdk.Context, k Keeper) (sdk.Dec, sdk.Dec, sdk.Dec) {
	// fetch params
	//params := k.GetParams(ctx)
	//BlocksPerYear := params.BlocksPerYear
	// recalculate inflation rate
	totalSupply := k.sk.TotalTokens(ctx)
	//
	curBlkHeight := ctx.BlockHeight()
	fmt.Println("curBlkHeight:", curBlkHeight)
	// BlockProvision = 25 / 2**roundIndex
	BlkRewardInt64 := int64(Htdf2Satoshi * calcMiningReward(int(curBlkHeight)))
	BlockProvision := sdk.NewDec(BlkRewardInt64)
	fmt.Println("BlockProvision:", BlkRewardInt64)
	AnnualProvisionsDec := sdk.NewDec(int64(AnnualProvisions))
	// Inflation = AnnualProvisions / totalSupply
	Inflation := AnnualProvisionsDec.Quo(sdk.NewDecFromInt(totalSupply))
	return AnnualProvisionsDec, Inflation, BlockProvision
}

// Inflate every block, update inflation parameters once per hour
func BeginBlocker(ctx sdk.Context, k Keeper) {

	// fetch stored minter & params
	minter := k.GetMinter(ctx)
	//params := k.GetParams(ctx)

	// recalculate inflation rate
	var provisionAmt sdk.Dec
	minter.AnnualProvisions, minter.Inflation, provisionAmt = calcParams(ctx, k)

	k.SetMinter(ctx, minter)

	// mint coins, add to collected fees, update supply
	//fmt.Printf("AnnualProvisions: %s, Inflation: %s, provisionAmt: %s\n", minter.AnnualProvisions.String(), minter.Inflation.String(), provisionAmt.TruncateInt().String())
	mintedCoin := sdk.NewCoin(sdk.DefaultDenom, provisionAmt.TruncateInt())
	k.fck.AddCollectedFees(ctx, sdk.Coins{mintedCoin})
	k.sk.InflateSupply(ctx, mintedCoin.Amount)

}

// Inflate every block, update inflation parameters once per hour
// func BeginBlocker(ctx sdk.Context, k Keeper) {

// 	// fetch stored minter & params
// 	minter := k.GetMinter(ctx)
// 	params := k.GetParams(ctx)

// 	// recalculate inflation rate
// 	totalSupply := k.sk.TotalTokens(ctx)
// 	bondedRatio := k.sk.BondedRatio(ctx)
// 	minter.Inflation = minter.NextInflationRate(params, bondedRatio)
// 	minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalSupply)
// 	k.SetMinter(ctx, minter)

// 	// mint coins, add to collected fees, update supply
// 	mintedCoin := minter.BlockProvision(params)
// 	k.fck.AddCollectedFees(ctx, sdk.Coins{mintedCoin})
// 	k.sk.InflateSupply(ctx, mintedCoin.Amount)

// }
