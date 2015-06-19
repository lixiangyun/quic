package crypto

import "github.com/romain-jacotin/quic/protocol"

type AEAD_ChaCha20Poly1305 struct {
}

// NewAEAD_ChaCha20Poly1305 returns a *AEAD_ChaCha20 that implements AEAD interface
func NewAEAD_ChaCha20() AEAD {
	return new(AEAD_ChaCha20Poly1305)
}

// SetKey
func (this *AEAD_ChaCha20Poly1305) SetKey(key []byte) error {
	return nil
}

//SetNoncePrefix
func (this *AEAD_ChaCha20Poly1305) SetNoncePrefix(nonce []byte) error {
	return nil
}

// Open
func (this *AEAD_ChaCha20Poly1305) Open(sequencenumber protocol.QuicPacketSequenceNumber, cleartext, associateddata, ciphertext, tag []byte) error {
	return nil
}

// Seal
func (this *AEAD_ChaCha20Poly1305) Seal(sequencenumber protocol.QuicPacketSequenceNumber, ciphertext, tag, associateddata, cleartext []byte) error {
	return nil
}