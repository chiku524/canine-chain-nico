package params

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/x/tx/signing"
	gogoproto "github.com/cosmos/gogoproto/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	filetreetypes "github.com/jackalLabs/canine-chain/v5/x/filetree/types"
	notificationstypes "github.com/jackalLabs/canine-chain/v5/x/notifications/types"
	oracletypes "github.com/jackalLabs/canine-chain/v5/x/oracle/types"
	rnstypes "github.com/jackalLabs/canine-chain/v5/x/rns/types"
	storagetypes "github.com/jackalLabs/canine-chain/v5/x/storage/types"
)

type hasCreator interface {
	GetCreator() string
}

// registerJackalCustomSigners wires GetSigners for Jackal msgs that lack
// cosmos.msg.v1.signer proto options (required by SDK 0.50+ tx signing).
func registerJackalCustomSigners(opts *signing.Options, accountPrefix string) {
	ac := address.Bech32Codec{Bech32Prefix: accountPrefix}
	fn := func(msg proto.Message) ([][]byte, error) {
		var creator string
		if c, ok := msg.(hasCreator); ok {
			creator = c.GetCreator()
		} else {
			refl := msg.ProtoReflect()
			fd := refl.Descriptor().Fields().ByName("creator")
			if fd == nil {
				return nil, fmt.Errorf("message %T has no creator field", msg)
			}
			creator = refl.Get(fd).String()
		}
		bz, err := ac.StringToBytes(creator)
		if err != nil {
			return nil, err
		}
		return [][]byte{bz}, nil
	}

	msgs := []gogoproto.Message{
		// storage
		&storagetypes.MsgPostFile{},
		&storagetypes.MsgPostProof{},
		&storagetypes.MsgPostProofFor{},
		&storagetypes.MsgDeleteFile{},
		&storagetypes.MsgSetProviderIP{},
		&storagetypes.MsgSetProviderKeybase{},
		&storagetypes.MsgSetProviderTotalSpace{},
		&storagetypes.MsgInitProvider{},
		&storagetypes.MsgShutdownProvider{},
		&storagetypes.MsgBuyStorage{},
		&storagetypes.MsgAddClaimer{},
		&storagetypes.MsgRemoveClaimer{},
		&storagetypes.MsgRequestAttestationForm{},
		&storagetypes.MsgAttest{},
		&storagetypes.MsgRequestReportForm{},
		&storagetypes.MsgReport{},
		// filetree
		&filetreetypes.MsgPostFile{},
		&filetreetypes.MsgAddViewers{},
		&filetreetypes.MsgPostKey{},
		&filetreetypes.MsgDeleteFile{},
		&filetreetypes.MsgRemoveViewers{},
		&filetreetypes.MsgProvisionFileTree{},
		&filetreetypes.MsgAddEditors{},
		&filetreetypes.MsgRemoveEditors{},
		&filetreetypes.MsgResetEditors{},
		&filetreetypes.MsgResetViewers{},
		&filetreetypes.MsgChangeOwner{},
		// rns
		&rnstypes.MsgRegister{},
		&rnstypes.MsgRegisterName{},
		&rnstypes.MsgUpdate{},
		&rnstypes.MsgMakePrimary{},
		&rnstypes.MsgBid{},
		&rnstypes.MsgAcceptBid{},
		&rnstypes.MsgCancelBid{},
		&rnstypes.MsgList{},
		&rnstypes.MsgBuy{},
		&rnstypes.MsgDelist{},
		&rnstypes.MsgTransfer{},
		&rnstypes.MsgAddRecord{},
		&rnstypes.MsgDelRecord{},
		&rnstypes.MsgInit{},
		// oracle
		&oracletypes.MsgCreateFeed{},
		&oracletypes.MsgUpdateFeed{},
		// notifications
		&notificationstypes.MsgCreateNotification{},
		&notificationstypes.MsgDeleteNotification{},
		&notificationstypes.MsgBlockSenders{},
	}

	for _, m := range msgs {
		name := gogoproto.MessageName(m)
		if name == "" {
			panic(fmt.Sprintf("empty proto name for %T", m))
		}
		opts.DefineCustomGetSigners(protoreflect.FullName(name), fn)
	}
}
