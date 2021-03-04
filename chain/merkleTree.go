package chain

import (
	"bytes"
	"crypto/sha256"
)

// SetHMerkleRoot calculate the hash of the merkle tree
func (b *Block) SetHMerkleRoot() {
	NTransaction := b.CountTransaction()
	switch NTransaction {
	case 0:
		HMerkleRoot := sha256.Sum256([]byte{0})
		b.Header.HMerkleRoot = HMerkleRoot[:]
	case 1:
		final := hConcat(b.Transactions[0].Serialize(), b.Transactions[0].Serialize())
		MerkleRoot := sha256.Sum256(final)
		b.Header.HMerkleRoot = MerkleRoot[:]
	case 2:
		final := hConcat(b.Transactions[0].Serialize(), b.Transactions[1].Serialize())
		MerkleRoot := sha256.Sum256(final)
		b.Header.HMerkleRoot = MerkleRoot[:]
	case 3:
		S1 := hConcat(b.Transactions[0].Serialize(), b.Transactions[1].Serialize())
		S2 := hConcat(b.Transactions[2].Serialize(), b.Transactions[2].Serialize())
		final := hConcat(S1, S2)
		MerkleRoot := sha256.Sum256(final)
		b.Header.HMerkleRoot = MerkleRoot[:]
	case 4:
		S1 := hConcat(b.Transactions[0].Serialize(), b.Transactions[1].Serialize())
		S2 := hConcat(b.Transactions[2].Serialize(), b.Transactions[3].Serialize())
		final := hConcat(S1, S2)
		MerkleRoot := sha256.Sum256(final)
		b.Header.HMerkleRoot = MerkleRoot[:]
	case 5:
		S1 := hConcat(b.Transactions[0].Serialize(), b.Transactions[1].Serialize())
		S2 := hConcat(b.Transactions[2].Serialize(), b.Transactions[3].Serialize())
		S3 := hConcat(b.Transactions[4].Serialize(), b.Transactions[4].Serialize())
		C1 := hConcat(S1, S2)
		C2 := hConcat(S3, S3)
		final := hConcat(C1, C2)
		MerkleRoot := sha256.Sum256(final)
		b.Header.HMerkleRoot = MerkleRoot[:]
	case 6:
		S1 := hConcat(b.Transactions[0].Serialize(), b.Transactions[1].Serialize())
		S2 := hConcat(b.Transactions[2].Serialize(), b.Transactions[3].ID)
		S3 := hConcat(b.Transactions[4].Serialize(), b.Transactions[5].Serialize())
		C1 := hConcat(S1, S2)
		C2 := hConcat(S3, S3)
		final := hConcat(C1, C2)
		MerkleRoot := sha256.Sum256(final)
		b.Header.HMerkleRoot = MerkleRoot[:]
	case 7:
		S1 := hConcat(b.Transactions[0].Serialize(), b.Transactions[1].Serialize())
		S2 := hConcat(b.Transactions[2].Serialize(), b.Transactions[3].Serialize())
		S3 := hConcat(b.Transactions[4].Serialize(), b.Transactions[5].Serialize())
		S4 := hConcat(b.Transactions[6].Serialize(), b.Transactions[6].Serialize())
		C1 := hConcat(S1, S2)
		C2 := hConcat(S3, S4)
		final := hConcat(C1, C2)
		MerkleRoot := sha256.Sum256(final)
		b.Header.HMerkleRoot = MerkleRoot[:]

	}
}

func hConcat(a, b []byte) []byte {
	concatenated := sha256.Sum256(bytes.Join(
		[][]byte{
			a,
			b,
		},
		[]byte{},
	))
	return concatenated[:]
}
