package mint

import (
	"fmt"

	sdk "github.com/orientwalt/htdf/types"
)

// junying-todo, 2019-07-17
//	BlksPerRound = 100
//	rewards+commission+community-pool
//	hscli query distr rewards htdf1zulqmaqlsgrgmagenaqf02p8kfgsuqkdwgwj80
//	121793749706.0satoshi
//	* 4 = 487174998824
//	not true becasue proper get more rewards,that's, different rewards on every node.
//  hscli query distr commission cosmosvaloper1lwjmdnks33xwnmfayc64ycprww49n33mtm92ne
// 	hscli query distr community-pool

// junying-todo, 2019-07-15
// single node: 88.2 for delegators, 11.8 for validator(commission)
// commission is validating fee
// commission rate changes?
func calcParams(ctx sdk.Context, k Keeper) (sdk.Dec, sdk.Dec, sdk.Dec) {
	// fetch params
	totalSupply := k.sk.TotalTokens(ctx)
	fmt.Println("totalSupply", totalSupply)
	// block index
	curBlkHeight := ctx.BlockHeight()
	fmt.Println("curBlkHeight:", curBlkHeight)

	// check terminate condition, junying-todo, 2019-12-05
	if totalSupply.GT(sdk.NewInt(TotalLiquidAsSatoshi)) { // || curBlkHeight > TotalMinableBlks {
		return sdk.NewDec(0), sdk.NewDec(0), sdk.NewDec(0)
	}

	AnnualProvisionsDec := sdk.NewDec(int64(AnnualProvisions))
	// Inflation = AnnualProvisions / totalSupply
	Inflation := AnnualProvisionsDec.Quo(sdk.NewDecFromInt(totalSupply))

	// sine params
	curAmplitude := k.sk.Amplitude(ctx)
	curCycle := k.sk.Cycle(ctx)
	curLastIndex := k.sk.LastIndex(ctx)
	// check if it's time for new cycle
	if curBlkHeight >= (curLastIndex + curCycle) {
		k.sk.SetAmplitude(ctx, randomAmplitude())
		k.sk.SetCycle(ctx, randomCycle())
		k.sk.SetLastIndex(ctx, curBlkHeight)
	}

	BlockReward := calcRewardAsSatoshi(curAmplitude, curCycle, curBlkHeight-curLastIndex)
	BlockProvision := sdk.NewDec(BlockReward)
	fmt.Println("BlockProvision:", BlockReward)
	fmt.Println("curAmplitude:", curAmplitude)
	fmt.Println("curCycle:", curCycle)
	fmt.Println("curLastIndex:", curLastIndex)

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
