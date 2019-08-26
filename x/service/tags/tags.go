package tags

import (
	sdk "github.com/orientwalt/htdf/types"
)

var (
	ActionSvcCallTimeOut = "service-call-expiration"

	Action = sdk.TagAction

	Provider   = "provider"
	Consumer   = "consumer"
	RequestID  = "request-id"
	ServiceFee = "service-fee"
	SlashCoins = "service-slash-coins"
)
