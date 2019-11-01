package crypto

// PrivateKey defines the interface of privatekey.
type PrivateKey interface {
	Public() PublicKey
	Sign([]byte) (Signature, error)
	SecretBytes() []byte
}

// PublicKey defines the interface of publickey.
type PublicKey interface {
	Bytes() []byte
}

// Signature defines the interface of signature.
type Signature interface {
	Verify([]byte, PublicKey) bool
	SetBytes([]byte)
	Bytes() []byte
}