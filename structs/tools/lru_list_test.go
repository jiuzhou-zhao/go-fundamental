package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TInvWitnessFlag = 1 << 30
	THashSize       = 32
)

// TInvType represents the allowed types of inventory vectors.  See InvVect.
type TInvType uint32

// These constants define the various supported inventory vector types.
const (
	TInvTypeError                TInvType = 0
	TInvTypeTx                   TInvType = 1
	TInvTypeBlock                TInvType = 2
	TInvTypeFilteredBlock        TInvType = 3
	TInvTypeWitnessBlock         TInvType = TInvTypeBlock | TInvWitnessFlag
	TInvTypeWitnessTx            TInvType = TInvTypeTx | TInvWitnessFlag
	TInvTypeFilteredWitnessBlock TInvType = TInvTypeFilteredBlock | TInvWitnessFlag
)

type THash [THashSize]byte

func (hash *THash) SetBytes(newHash []byte) error {
	copy(hash[:], newHash)
	return nil
}

type TInvVect struct {
	Type TInvType // Type of data
	Hash THash    // Hash of the data
}

func NewHash(newHash []byte) (*THash, error) {
	var sh THash
	err := sh.SetBytes(newHash)
	if err != nil {
		return nil, err
	}
	return &sh, err
}

func TestLruLis(t *testing.T) {
	l := NewLruList(3)

	h1, _ := NewHash([]byte{0x01})
	h2, _ := NewHash([]byte{0x01, 0x02})
	h3, _ := NewHash([]byte{0x01, 0x02, 0x03})

	l.Add(TInvVect{
		Type: TInvTypeTx,
		Hash: *h1,
	}, 10)
	l.Add(TInvVect{
		Type: TInvTypeTx,
		Hash: *h2,
	}, 10)

	assert.True(t, l.Exists(TInvVect{
		Type: TInvTypeTx,
		Hash: *h1,
	}))
	assert.True(t, l.Exists(TInvVect{
		Type: TInvTypeTx,
		Hash: *h2,
	}))
	assert.False(t, l.Exists(TInvVect{
		Type: TInvTypeTx,
		Hash: *h3,
	}))

	l.Add(TInvVect{
		Type: TInvTypeBlock,
		Hash: *h3,
	}, 10)
	assert.True(t, l.Exists(TInvVect{
		Type: TInvTypeBlock,
		Hash: *h3,
	}))

	l.Add(TInvVect{
		Type: TInvTypeFilteredBlock,
		Hash: *h3,
	}, 10)
	assert.True(t, l.Exists(TInvVect{
		Type: TInvTypeBlock,
		Hash: *h3,
	}))
	assert.True(t, l.Exists(TInvVect{
		Type: TInvTypeFilteredBlock,
		Hash: *h3,
	}))
	assert.False(t, l.Exists(TInvVect{
		Type: TInvTypeTx,
		Hash: *h1,
	}))
}
