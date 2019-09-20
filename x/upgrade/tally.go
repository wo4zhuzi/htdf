package upgrade

import (
	"fmt"

	sdk "github.com/orientwalt/htdf/types"
)

func tally(ctx sdk.Context, versionProtocol uint64, k Keeper, threshold sdk.Dec) (passes bool) {

	totalVotingPower := sdk.ZeroInt()
	signalsVotingPower := sdk.ZeroInt()
	k.sk.IterateBondedValidatorsByPower(ctx, func(index int64, validator sdk.Validator) (stop bool) {
		totalVotingPower = totalVotingPower.AddRaw(validator.GetTendermintPower())
		valAcc := validator.GetConsAddr().String()
		if ok := k.GetSignal(ctx, versionProtocol, valAcc); ok {
			signalsVotingPower = signalsVotingPower.AddRaw(validator.GetTendermintPower())
		}
		fmt.Print("7777777777777777777	", validator.GetTendermintPower(), "sersion protocol ", versionProtocol, "\n")
		fmt.Print("7777777777777777777	", valAcc, "\n")
		return false
	})

	ctx.Logger().Info("Tally Start", "SiganlsVotingPower", signalsVotingPower.String(),
		"TotalVotingPower", totalVotingPower.String(),
		"SiganlsVotingPower/TotalVotingPower", signalsVotingPower.Quo(totalVotingPower).String(),
		"Threshold", threshold.String())
	// If more than 95% of validator update , do switch
	fmt.Print("77777signalsVotingPower7777	", signalsVotingPower, "\n")
	fmt.Print("77777totalVotingPowerr7777	", totalVotingPower, "\n")
	fmt.Print("77777threshold7777	", threshold.RoundInt(), "\n")
	if signalsVotingPower.Quo(totalVotingPower).GT(threshold.RoundInt()) {
		return true
	}
	return true
}
