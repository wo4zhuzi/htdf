package cli

import (
	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/client/utils"
	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	authtxb "github.com/orientwalt/htdf/x/auth/client/txbuilder"
	"github.com/orientwalt/htdf/x/slashing"
	hscorecli "github.com/orientwalt/htdf/x/core/client/cli"
	"github.com/spf13/cobra"
)

// GetCmdUnjail implements the create unjail validator command.
func GetCmdUnjail(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "unjail [keyaddr]",
		Short: "unjail validator previously jailed for downtime",
		Long: `unjail a jailed validator:

$ hscli tx slashing unjail [keyaddr]
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			validatorAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := slashing.NewMsgUnjail(sdk.ValAddress(validatorAddr))
			return hscorecli.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg}, validatorAddr)
		},
	}
}
