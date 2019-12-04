package gov

import (
	"fmt"
	"strconv"

	sdk "github.com/orientwalt/htdf/types"
	"github.com/orientwalt/htdf/x/gov/tags"
)

// Handle all "gov" type messages.
func NewHandler(keeper Keeper) sdk.Handler {

	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgDeposit:
			fmt.Println("++++++++++++++++++++++++", msg)
			return handleMsgDeposit(ctx, keeper, msg)
		case MsgSubmitProposal:
			fmt.Println("++++++++++++++++++++++++", msg)
			return handleMsgSubmitProposal(ctx, keeper, msg)
		case MsgSubmitSoftwareUpgradeProposal:
			fmt.Println("++++++++++++++++++++++++", msg)
			return handleMsgSubmitSoftwareUpgradeProposal(ctx, keeper, msg)
		case MsgVote:
			return handleMsgVote(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized gov msg type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSubmitProposal(ctx sdk.Context, keeper Keeper, msg MsgSubmitProposal) sdk.Result {
	var content ProposalContent
	switch msg.ProposalType {
	case ProposalTypeText:
		content = NewTextProposal(msg.Title, msg.Description)
	default:
		return ErrInvalidProposalType(keeper.codespace, msg.ProposalType).Result()
	}
	proposal, err := keeper.SubmitProposal(ctx, content)
	if err != nil {
		return err.Result()
	}
	proposalID := proposal.GetProposalID()
	proposalIDStr := fmt.Sprintf("%d", proposalID)

	err, votingStarted := keeper.AddDeposit(ctx, proposalID, msg.Proposer, msg.InitialDeposit)
	if err != nil {
		return err.Result()
	}

	resTags := sdk.NewTags(
		tags.Proposer, []byte(msg.Proposer.String()),
		tags.ProposalID, proposalIDStr,
	)

	if votingStarted {
		resTags = resTags.AppendTag(tags.VotingPeriodStart, proposalIDStr)
	}

	return sdk.Result{
		Data: keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID),
		Tags: resTags,
	}
	fmt.Println("..........handleMsgSubmitSoftwareUpgradeProposal..........")
	return sdk.Result{}
}

//Submit upgrade software proposal
func handleMsgSubmitSoftwareUpgradeProposal(ctx sdk.Context, keeper Keeper, msg MsgSubmitSoftwareUpgradeProposal) sdk.Result {
	fmt.Println("1--------------handleMsgSubmitSoftwareUpgradeProposal---------------------", msg.Version)
	if !keeper.protocolKeeper.IsValidVersion(ctx, msg.Version) {

		return ErrCodeInvalidVersion(keeper.codespace, msg.Version).Result()
	}
	fmt.Println("3--------------handleMsgSubmitSoftwareUpgradeProposal---------------------")
	if uint64(ctx.BlockHeight()) > msg.SwitchHeight {
		return ErrCodeInvalidSwitchHeight(keeper.codespace, uint64(ctx.BlockHeight()), msg.SwitchHeight).Result()
	}
	fmt.Println("4--------------handleMsgSubmitSoftwareUpgradeProposal---------------------")
	_, found := keeper.guardianKeeper.GetProfiler(ctx, msg.Proposer)
	if !found {
		return ErrNotProfiler(keeper.codespace, msg.Proposer).Result()
	}
	fmt.Println("5--------------handleMsgSubmitSoftwareUpgradeProposal---------------------")
	if _, ok := keeper.protocolKeeper.GetUpgradeConfig(ctx); ok {
		return ErrSwitchPeriodInProcess(keeper.codespace).Result()
	}
	fmt.Println("6--------------handleMsgSubmitSoftwareUpgradeProposal---------------------")
	proposal := keeper.NewSoftwareUpgradeProposal(ctx, msg)
	fmt.Println("7--------------handleMsgSubmitSoftwareUpgradeProposal---------------------")
	err, votingStarted := keeper.AddInitialDeposit(ctx, proposal, msg.Proposer, msg.InitialDeposit)
	fmt.Println("8--------------handleMsgSubmitSoftwareUpgradeProposal---------------------", proposal.GetProposalID())
	if err != nil {
		return err.Result()
	}
	proposalIDBytes := strconv.FormatUint(proposal.GetProposalID(), 10)

	resTags := sdk.NewTags(
		tags.Proposer, []byte(msg.Proposer.String()),
		tags.ProposalID, proposalIDBytes,
	)

	if votingStarted {
		resTags = resTags.AppendTag(tags.VotingPeriodStart, proposalIDBytes)
	}
	fmt.Println("END--------------handleMsgSubmitSoftwareUpgradeProposal---------------------")
	// keeper.AddProposalNum(ctx, proposal)
	return sdk.Result{
		Data: []byte(proposalIDBytes),
		Tags: resTags,
	}
}

func handleMsgDeposit(ctx sdk.Context, keeper Keeper, msg MsgDeposit) sdk.Result {
	err, votingStarted := keeper.AddDeposit(ctx, msg.ProposalID, msg.Depositor, msg.Amount)
	if err != nil {
		return err.Result()
	}

	proposalIDStr := fmt.Sprintf("%d", msg.ProposalID)
	resTags := sdk.NewTags(
		tags.Depositor, []byte(msg.Depositor.String()),
		tags.ProposalID, proposalIDStr,
	)

	if votingStarted {
		resTags = resTags.AppendTag(tags.VotingPeriodStart, proposalIDStr)
	}

	return sdk.Result{
		Tags: resTags,
	}
}

func handleMsgVote(ctx sdk.Context, keeper Keeper, msg MsgVote) sdk.Result {
	err := keeper.AddVote(ctx, msg.ProposalID, msg.Voter, msg.Option)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: sdk.NewTags(
			tags.Voter, msg.Voter.String(),
			tags.ProposalID, fmt.Sprintf("%d", msg.ProposalID),
		),
	}
}
