// nolint
package tags

import (
	sdk "github.com/orientwalt/htdf/types"
)

var (
	ActionCompleteUnbonding     = "complete-unbonding"
	ActionCompleteRedelegation  = "complete-redelegation"
	ActionCompleteAuthorization = "complete-authorization"

	Action       = sdk.TagAction
	SrcValidator = sdk.TagSrcValidator
	DstValidator = sdk.TagDstValidator
	Delegator    = sdk.TagDelegator
	Moniker      = "moniker"
	Identity     = "identity"
	EndTime      = "end-time"
)
