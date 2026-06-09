package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestOverflow_Finding3 documents Finding #3: naive FileSize*MaxProofs can wrap in Go int64.
func TestOverflow_Finding3(t *testing.T) {
	fileSize := int64(1 << 40)
	maxProofs := int64(1 << 25)
	wrapped := fileSize * maxProofs
	t.Logf("raw multiply wraps: %d * %d = %d", fileSize, maxProofs, wrapped)
	require.Less(t, wrapped, fileSize, "wrapped product must be smaller than file size")

	_, err := mulStorageCharge(fileSize, maxProofs)
	require.Error(t, err, "keeper must reject overflowed storage charges")

	fileSize = 1 << 50
	maxProofs = 1 << 14
	wrapped = fileSize * maxProofs
	t.Logf("raw multiply wraps to small value: %d * %d = %d", fileSize, maxProofs, wrapped)
	require.Less(t, wrapped, fileSize)

	_, err = mulStorageCharge(fileSize, maxProofs)
	require.Error(t, err)
}
