package signs

import (
	"fmt"
	"github.com/orientwalt/htdf/client/context"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
)

//VerifySignature return an Tx validate
func VerifySignature(cliCtx context.CLIContext, chainID string, stdTx auth.StdTx) bool {

	fmt.Println("Signers:")

	signers := stdTx.GetSigners()
	for i, signer := range signers {
		fmt.Printf(" %v: %v\n", i, signer.String())
	}

	success := true
	sigs := stdTx.GetSignatures()

	fmt.Println("")
	fmt.Println("Signatures:")

	if len(sigs) != len(signers) {
		success = false
	}

	for i, sig := range sigs {
		sigAddr := sdk.AccAddress(sig.Address())
		sigSanity := "OK"

		if i >= len(signers) || !sigAddr.Equals(signers[i]) {
			sigSanity = "ERROR: signature does not match its respective signer"
			success = false
		}

		// Validate the actual signature over the transaction bytes since we can
		// reach out to a full node to query accounts.
		if success {
			acc, err := cliCtx.GetAccount(sigAddr)
			if err != nil {
				fmt.Printf("failed to get account: %s\n", sigAddr)
				return false
			}

			sigBytes := auth.StdSignBytes(
				chainID, acc.GetAccountNumber(), acc.GetSequence(),
				stdTx.Fee, stdTx.GetMsgs(), stdTx.GetMemo(),
			)

			if ok := sig.VerifyBytes(sigBytes, sig.Signature); !ok {
				sigSanity = "ERROR: signature invalid"
				success = false
			}
		}

		fmt.Printf(" %v: %v\t[%s]\n", i, sigAddr.String(), sigSanity)
	}

	fmt.Println("")
	return success
}
