// nolint
package tags

import (
	sdk "github.com/orientwalt/htdf/types"
)

// Distribution tx tags
var (
	Rewards    = "rewards"
	Commission = "commission"

	Validator = sdk.TagSrcValidator
	Delegator = sdk.TagDelegator
)
