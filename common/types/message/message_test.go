package message

import (
	"encoding/hex"
	"testing"

	"github.com/scroll-tech/go-ethereum/common"
	"github.com/scroll-tech/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestAuthMessageSignAndVerify(t *testing.T) {
	privkey, err := crypto.GenerateKey()
	assert.NoError(t, err)

	authMsg := &AuthMsg{
		Identity: &Identity{
			Name:    "testName",
			Version: "testVersion",
			Token:   "testToken",
		},
	}
	assert.NoError(t, authMsg.SignWithKey(privkey))

	// Check public key.
	pk, err := authMsg.PublicKey()
	assert.NoError(t, err)
	assert.Equal(t, common.Bytes2Hex(crypto.CompressPubkey(&privkey.PublicKey)), pk)

	ok, err := authMsg.Verify()
	assert.NoError(t, err)
	assert.Equal(t, true, ok)

	// Check public key is ok.
	pub, err := authMsg.PublicKey()
	assert.NoError(t, err)
	pubkey := crypto.CompressPubkey(&privkey.PublicKey)
	assert.Equal(t, pub, common.Bytes2Hex(pubkey))
}

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken()
	assert.NoError(t, err)
	assert.Equal(t, 32, len(token))
}

func TestIdentityHash(t *testing.T) {
	identity := &Identity{
		Name:       "testName",
		ProverType: ProofTypeChunk,
		Version:    "testVersion",
		Token:      "testToken",
	}
	hash, err := identity.Hash()
	assert.NoError(t, err)

	expectedHash := "c0411a19531fb8c6133b2bae91f361c14e65f2d318aef72b83519e6061cad001"
	assert.Equal(t, expectedHash, hex.EncodeToString(hash))
}

func TestProofMessageSignVerifyPublicKey(t *testing.T) {
	privkey, err := crypto.GenerateKey()
	assert.NoError(t, err)

	proofMsg := &ProofMsg{
		ProofDetail: &ProofDetail{
			ID:     "testID",
			Type:   ProofTypeChunk,
			Status: StatusOk,
			ChunkProof: &ChunkProof{
				StorageTrace: []byte("testStorageTrace"),
				Protocol:     []byte("testProtocol"),
				Proof:        []byte("testProof"),
				Instances:    []byte("testInstance"),
				Vk:           []byte("testVk"),
			},
			Error: "testError",
		},
	}
	assert.NoError(t, proofMsg.Sign(privkey))

	// Test when publicKey is not set.
	ok, err := proofMsg.Verify()
	assert.NoError(t, err)
	assert.Equal(t, true, ok)

	// Test when publicKey is already set.
	ok, err = proofMsg.Verify()
	assert.NoError(t, err)
	assert.Equal(t, true, ok)
}

func TestProofDetailHash(t *testing.T) {
	proofDetail := &ProofDetail{
		ID:     "testID",
		Type:   ProofTypeChunk,
		Status: StatusOk,
		ChunkProof: &ChunkProof{
			StorageTrace: []byte("testStorageTrace"),
			Protocol:     []byte("testProtocol"),
			Proof:        []byte("testProof"),
			Instances:    []byte("testInstance"),
			Vk:           []byte("testVk"),
		},
		Error: "testError",
	}
	hash, err := proofDetail.Hash()
	assert.NoError(t, err)
	expectedHash := "4165f5ab3399001002a5b8e4062914249a2deb72f6133d647b586f53e236802d"
	assert.Equal(t, expectedHash, hex.EncodeToString(hash))
}

func TestProveTypeString(t *testing.T) {
	proofTypeChunk := ProofType(1)
	assert.Equal(t, "proof type chunk", proofTypeChunk.String())

	proofTypeBatch := ProofType(2)
	assert.Equal(t, "proof type batch", proofTypeBatch.String())

	illegalProof := ProofType(3)
	assert.Equal(t, "illegal proof type", illegalProof.String())
}

func TestProofMsgPublicKey(t *testing.T) {
	privkey, err := crypto.GenerateKey()
	assert.NoError(t, err)

	proofMsg := &ProofMsg{
		ProofDetail: &ProofDetail{
			ID:     "testID",
			Type:   ProofTypeChunk,
			Status: StatusOk,
			ChunkProof: &ChunkProof{
				StorageTrace: []byte("testStorageTrace"),
				Protocol:     []byte("testProtocol"),
				Proof:        []byte("testProof"),
				Instances:    []byte("testInstance"),
				Vk:           []byte("testVk"),
			},
			Error: "testError",
		},
	}
	assert.NoError(t, proofMsg.Sign(privkey))

	// Test when publicKey is not set.
	pk, err := proofMsg.PublicKey()
	assert.NoError(t, err)
	assert.Equal(t, common.Bytes2Hex(crypto.CompressPubkey(&privkey.PublicKey)), pk)

	// Test when publicKey is already set.
	proofMsg.publicKey = common.Bytes2Hex(crypto.CompressPubkey(&privkey.PublicKey))
	pk, err = proofMsg.PublicKey()
	assert.NoError(t, err)
	assert.Equal(t, common.Bytes2Hex(crypto.CompressPubkey(&privkey.PublicKey)), pk)
}
