package keeper

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/jackalLabs/canine-chain/v5/x/rns/types"
)

func (k msgServer) List(goCtx context.Context, msg *types.MsgList) (*types.MsgListResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	mname := strings.ToLower(msg.Name)

	_, found := k.GetForsale(ctx, mname)

	if found {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Name already listed.")
	}

	n, tld, err := GetNameAndTLD(mname)
	if err != nil {
		return nil, err
	}

	name, nfound := k.GetNames(ctx, n, tld)

	if !nfound {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Name does not exist or has expired.")
	}

	if name.Value != msg.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "You do not own this name.")
	}

	blockHeight := ctx.BlockHeight()

	if name.Locked > blockHeight {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "cannot transfer free name")
	}

	if blockHeight > name.Expires {
		return nil, errorsmod.Wrap(sdkerrors.ErrNotFound, "Name does not exist or has expired.")
	}

	newsale := types.Forsale{
		Name:  mname,
		Price: msg.Price.String(),
		Owner: msg.Creator,
	}

	k.SetForsale(ctx, newsale)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeJackalMessage,
			sdk.NewAttribute(types.AttributeKeySigner, msg.Creator),
		),
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventSetSale,
			sdk.NewAttribute(types.AttributeKeySigner, msg.Creator),
		),
	)

	return &types.MsgListResponse{}, nil
}
