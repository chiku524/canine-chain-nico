package app

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	storagetypes "github.com/jackalLabs/canine-chain/v5/x/storage/types"
	"github.com/stretchr/testify/require"
)

type mockTx struct {
	msgs []sdk.Msg
}

func (m mockTx) GetMsgs() []sdk.Msg            { return m.msgs }
func (m mockTx) GetMsgsV2() ([]sdk.Msg, error) { return m.msgs, nil }
func (m mockTx) ValidateBasic() error          { return nil }

func TestIsFreeStorageTx(t *testing.T) {
	t.Parallel()

	require.True(t, isFreeStorageTx(mockTx{msgs: []sdk.Msg{
		&storagetypes.MsgPostProof{},
	}}))
	require.True(t, isFreeStorageTx(mockTx{msgs: []sdk.Msg{
		&storagetypes.MsgPostProof{},
		&storagetypes.MsgAttest{},
	}}))
	require.False(t, isFreeStorageTx(mockTx{msgs: []sdk.Msg{
		&storagetypes.MsgPostProof{},
		&storagetypes.MsgBuyStorage{},
	}}))
	require.False(t, isFreeStorageTx(mockTx{msgs: []sdk.Msg{}}))
	require.False(t, isFreeStorageTx(mockTx{msgs: []sdk.Msg{
		&storagetypes.MsgPostProofFor{},
	}}))
}
