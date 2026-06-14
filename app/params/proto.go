package params

import (
	"github.com/cosmos/gogoproto/proto"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/tx/signing"
)

// MakeEncodingConfig creates an EncodingConfig using sdk.GetConfig() bech32 prefixes.
func MakeEncodingConfig() EncodingConfig {
	return MakeEncodingConfigWithBech32(
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		sdk.GetConfig().GetBech32ValidatorAddrPrefix(),
	)
}

// MakeEncodingConfigWithBech32 creates an EncodingConfig with explicit bech32 prefixes.
func MakeEncodingConfigWithBech32(accountPrefix, validatorPrefix string) EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry, err := types.NewInterfaceRegistryWithOptions(types.InterfaceRegistryOptions{
		ProtoFiles: proto.HybridResolver,
		SigningOptions: signing.Options{
			AddressCodec: address.Bech32Codec{
				Bech32Prefix: accountPrefix,
			},
			ValidatorAddressCodec: address.Bech32Codec{
				Bech32Prefix: validatorPrefix,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(marshaler, tx.DefaultSignModes)

	return EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
}
