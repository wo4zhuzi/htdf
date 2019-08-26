package cli

import (
	"fmt"
	"strings"

	"github.com/orientwalt/htdf/client"
	"github.com/orientwalt/htdf/client/context"
	"github.com/orientwalt/htdf/client/utils"
	"github.com/orientwalt/htdf/codec"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/auth"
	authtxb "github.com/orientwalt/htdf/x/auth/client/txbuilder"
	"github.com/orientwalt/htdf/x/staking"
	hscorecli "github.com/orientwalt/htdf/x/core/client/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdCreateValidator implements the create validator command handler.
// TODO: Add full description
func GetCmdCreateValidator(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validator [accaddr]",
		Short: "create new validator initialized with a self-delegation to it",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			validatorAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			txBldr, msg, err := BuildCreateValidatorMsg(cliCtx, txBldr, validatorAddr)
			if err != nil {
				return err
			}

			return hscorecli.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg}, validatorAddr)
		},
	}

	cmd.Flags().AddFlagSet(FsPk)
	cmd.Flags().AddFlagSet(FsAmount)
	cmd.Flags().AddFlagSet(fsDescriptionCreate)
	cmd.Flags().AddFlagSet(FsCommissionCreate)
	cmd.Flags().AddFlagSet(FsMinSelfDelegation)
	cmd.Flags().AddFlagSet(fsDelegator)
	cmd.Flags().String(FlagIP, "", fmt.Sprintf("The node's public IP. It takes effect only when used in combination with --%s", client.FlagGenerateOnly))
	cmd.Flags().String(FlagNodeID, "", "The node's ID")

	cmd.MarkFlagRequired(FlagAmount)
	cmd.MarkFlagRequired(FlagPubKey)
	cmd.MarkFlagRequired(FlagMoniker)

	return cmd
}

// GetCmdEditValidator implements the create edit validator command.
// TODO: add full description
func GetCmdEditValidator(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-validator [accaddr]",
		Short: "edit an existing validator account",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			validatorAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			description := staking.Description{
				Moniker:  viper.GetString(FlagMoniker),
				Identity: viper.GetString(FlagIdentity),
				Website:  viper.GetString(FlagWebsite),
				Details:  viper.GetString(FlagDetails),
			}

			var newRate *sdk.Dec

			commissionRate := viper.GetString(FlagCommissionRate)
			if commissionRate != "" {
				rate, err := sdk.NewDecFromStr(commissionRate)
				if err != nil {
					return fmt.Errorf("invalid new commission rate: %v", err)
				}

				newRate = &rate
			}

			var newMinSelfDelegation *sdk.Int

			minSelfDelegationString := viper.GetString(FlagMinSelfDelegation)
			if minSelfDelegationString != "" {
				msb, ok := sdk.NewIntFromString(minSelfDelegationString)
				if !ok {
					return fmt.Errorf(staking.ErrMinSelfDelegationInvalid(staking.DefaultCodespace).Error())
				}
				newMinSelfDelegation = &msb
			}

			msg := staking.NewMsgEditValidator(sdk.ValAddress(validatorAddr), description, newRate, newMinSelfDelegation)

			// build and sign the transaction, then broadcast to Tendermint
			return hscorecli.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg}, validatorAddr)
		},
	}

	cmd.Flags().AddFlagSet(fsDescriptionEdit)
	cmd.Flags().AddFlagSet(fsCommissionUpdate)

	return cmd
}

// GetCmdDelegate implements the delegate command.
func GetCmdDelegate(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "delegate [delegator-addr] [validator-addr] [amount]",
		Args:  cobra.ExactArgs(3),
		Short: "delegate liquid tokens to a validator",
		Long: strings.TrimSpace(`Delegate an amount of liquid coins to a validator from your wallet:

$ hscli tx staking delegate cosmos1tq7zajghkxct4al0yf44ua9rjwnw06vdusflk4 \
				   			cosmosvaloper1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm \
				   			1000stake
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			delAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			valAddr, err := sdk.ValAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoin(args[2])
			if err != nil {
				return err
			}

			msg := staking.NewMsgDelegate(delAddr, valAddr, amount)
			return hscorecli.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg}, delAddr)
		},
	}
}

// GetCmdRedelegate the begin redelegation command.
func GetCmdRedelegate(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "redelegate [delegator-addr] [src-validator-accaddr] [dst-validator-accaddr] [amount]",
		Short: "redelegate illiquid tokens from one validator to another",
		Args:  cobra.ExactArgs(4),
		Long: strings.TrimSpace(`Redelegate an amount of illiquid staking tokens from one validator to another:

$ hscli tx staking redelegate cosmos1tq7zajghkxct4al0yf44ua9rjwnw06vdusflk4 \
							  cosmosvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj \
							  cosmosvaloper1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm \
							  100stake
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			delAddr := cliCtx.GetFromAddress()
			valSrcAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			valDstAddr, err := sdk.ValAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoin(args[2])
			if err != nil {
				return err
			}

			msg := staking.NewMsgBeginRedelegate(delAddr, valSrcAddr, valDstAddr, amount)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}
}

// GetCmdUnbond implements the unbond validator command.
func GetCmdUnbond(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "unbond [delegator-addr] [validator-addr] [amount]",
		Short: "unbond shares from a validator",
		Args:  cobra.ExactArgs(3),
		Long: strings.TrimSpace(`Unbond an amount of bonded shares from a validator:

$ hscli tx staking unbond cosmos1tq7zajghkxct4al0yf44ua9rjwnw06vdusflk4 \
						  cosmosvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj \
						  100stake
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			delAddr := cliCtx.GetFromAddress()
			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoin(args[1])
			if err != nil {
				return err
			}

			msg := staking.NewMsgUndelegate(delAddr, valAddr, amount)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}
}

// BuildCreateValidatorMsg makes a new MsgCreateValidator.
func BuildCreateValidatorMsg(cliCtx context.CLIContext, txBldr authtxb.TxBuilder, valAddr sdk.AccAddress) (authtxb.TxBuilder, sdk.Msg, error) {
	amounstStr := viper.GetString(FlagAmount)
	amount, err := sdk.ParseCoin(amounstStr)
	if err != nil {
		return txBldr, nil, err
	}

	pkStr := viper.GetString(FlagPubKey)
	pk, err := sdk.GetConsPubKeyBech32(pkStr)
	if err != nil {
		return txBldr, nil, err
	}

	description := staking.NewDescription(
		viper.GetString(FlagMoniker),
		viper.GetString(FlagIdentity),
		viper.GetString(FlagWebsite),
		viper.GetString(FlagDetails),
	)

	// get the initial validator commission parameters
	rateStr := viper.GetString(FlagCommissionRate)
	maxRateStr := viper.GetString(FlagCommissionMaxRate)
	maxChangeRateStr := viper.GetString(FlagCommissionMaxChangeRate)
	commissionMsg, err := buildCommissionMsg(rateStr, maxRateStr, maxChangeRateStr)
	if err != nil {
		return txBldr, nil, err
	}

	// get the initial validator min self delegation
	msbStr := viper.GetString(FlagMinSelfDelegation)
	minSelfDelegation, ok := sdk.NewIntFromString(msbStr)
	if !ok {
		return txBldr, nil, fmt.Errorf(staking.ErrMinSelfDelegationInvalid(staking.DefaultCodespace).Error())
	}

	msg := staking.NewMsgCreateValidator(
		sdk.ValAddress(valAddr), pk, amount, description, commissionMsg, minSelfDelegation,
	)

	if viper.GetBool(client.FlagGenerateOnly) {
		ip := viper.GetString(FlagIP)
		nodeID := viper.GetString(FlagNodeID)
		if nodeID != "" && ip != "" {
			txBldr = txBldr.WithMemo(fmt.Sprintf("%s@%s:26656", nodeID, ip))
		}
	}
	return txBldr, msg, nil
}
