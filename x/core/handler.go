package htdfservice

import (
	"fmt"

	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/bank"
)

type StateTransition interface {
	Transition(ctx sdk.Context, msg sdk.Msg) sdk.Result
}

func NewHandler(stateTransition StateTransition) sdk.Handler {
	return stateTransition.Transition
}

func NewClassicHandler(keeper bank.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSendFrom:
			return HandleMsgSendFrom(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized htdfservice Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to sendfrom
func HandleMsgSendFrom(ctx sdk.Context, keeper bank.Keeper, msg MsgSendFrom) sdk.Result {
	if !keeper.GetSendEnabled(ctx) {
		return bank.ErrSendDisabled(keeper.Codespace()).Result()
	}
	tags, err := keeper.SendCoins(ctx, msg.From, msg.To, msg.Amount)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}
