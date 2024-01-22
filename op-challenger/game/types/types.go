package types

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-service/sources/batching"
	"github.com/ethereum/go-ethereum/common"
)

type GameStatus uint8

const (
	GameStatusInProgress GameStatus = iota
	GameStatusChallengerWon
	GameStatusDefenderWon
)

// String returns the string representation of the game status.
func (s GameStatus) String() string {
	switch s {
	case GameStatusInProgress:
		return "In Progress"
	case GameStatusChallengerWon:
		return "Challenger Won"
	case GameStatusDefenderWon:
		return "Defender Won"
	default:
		return "Unknown"
	}
}

// GameStatusFromUint8 returns a game status from the uint8 representation.
func GameStatusFromUint8(i uint8) (GameStatus, error) {
	if i > 2 {
		return GameStatus(i), fmt.Errorf("invalid game status: %d", i)
	}
	return GameStatus(i), nil
}

type GameMetadata struct {
	GameType  uint8
	Timestamp uint64
	Proxy     common.Address
}

type LargePreimageIdent struct {
	Claimant common.Address
	UUID     *big.Int
}

type LargePreimageMetaData struct {
	LargePreimageIdent

	// Timestamp is the time at which the proposal first became fully available.
	// 0 when not all data is available yet
	Timestamp       uint64
	PartOffset      uint32
	ClaimedSize     uint32
	BlocksProcessed uint32
	BytesProcessed  uint32
	Countered       bool
}

// ShouldVerify returns true if the preimage upload is complete and has not yet been countered.
// Note that the challenge period for the preimage may have expired but the image not yet been finalized.
func (m LargePreimageMetaData) ShouldVerify() bool {
	return m.Timestamp > 0 && !m.Countered
}

const LeafSize = 136

// Leaf is the keccak state matrix added to the large preimage merkle tree.
type Leaf struct {
	// Input is the data absorbed for the block, exactly 136 bytes
	Input [LeafSize]byte
	// Index of the block in the absorption process
	Index *big.Int
	// StateCommitment is the hash of the internal state after absorbing the input.
	StateCommitment common.Hash
}

type LargePreimageOracle interface {
	Addr() common.Address
	GetActivePreimages(ctx context.Context, blockHash common.Hash) ([]LargePreimageMetaData, error)
	GetLeafBlocks(ctx context.Context, block batching.Block, ident LargePreimageIdent) ([]uint64, error)
	DecodeLeafData(data []byte) (*big.Int, []Leaf, error)
}
