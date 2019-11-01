package secp256k1

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/jfjun/wallet-tools/crypto"
	"github.com/jfjun/wallet-tools/utils"
	"io"
	"io/ioutil"
	"math/big"
	"os"


	"github.com/jfjun/crypto-tools/math"
	"github.com/jfjun/crypto-tools/crypto/secp256k1"
)

const (
	// SignatureSize represents the signature length
	SignatureSize = 65
)

var (
	// N is secp256k1 N
	N = S256().Params().N
	// halfN is N / 2
	halfN          = new(big.Int).Rsh(N, 1)
	emptySignature = Signature{}
)

type (
	// PrivateKey represents the ecdsa privatekey
	PrivateKey ecdsa.PrivateKey
	// PublicKey represents the ecdsa publickey
	PublicKey ecdsa.PublicKey
	// Signature represents the ecdsa_signcompact signature
	// data format [r - s - v]
	Signature [SignatureSize]byte
)

// S256 returns an instance of the secp256k1 curve
func S256() elliptic.Curve {
	return secp256k1.S256()
}

// GenerateKey returns a random PrivateKey
func GenerateKey() (*PrivateKey, error) {
	priv, err := ecdsa.GenerateKey(S256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return (*PrivateKey)(priv), err
}

// SecretBytes returns the actual bytes of ecdsa privatekey
func (priv *PrivateKey) SecretBytes() []byte {
	return math.PaddedBigBytes(priv.D, priv.Params().BitSize/8)
}

// Sign signs the hash and returns the signature
func (priv *PrivateKey) Sign(hash []byte) (s crypto.Signature, err error) {
	if len(hash) != 32 {
		return &emptySignature, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(hash))
	}
	secretKey := priv.SecretBytes()
	defer utils.ZeroMemory(secretKey)

	rawSig, err := secp256k1.Sign(hash, secretKey)

	sig := new(Signature)
	sig.compact(rawSig, false)

	return sig, err
}

// Public returns the public key corresponding to priv.
func (priv *PrivateKey) Public() crypto.PublicKey {
	return (*PublicKey)(&priv.PublicKey)
}

// SaveECDSA saves a private key to the given file
func (priv *PrivateKey) SaveECDSA(file string) error {
	ioutil.WriteFile(file, []byte(hex.EncodeToString(priv.SecretBytes())), 0600)
	return nil
}

// LoadECDSA loads a private key from the given file
func LoadECDSA(file string) (*PrivateKey, error) {
	buf := make([]byte, 64)
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	if _, err = io.ReadFull(fd, buf); err != nil {
		return nil, err
	}

	return HexToECDSA(string(buf))
}

// HexToECDSA parses a secp256k1 private key
func HexToECDSA(hexkey string) (*PrivateKey, error) {
	b, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, errors.New("invalid hex string")
	}
	if len(b) != 32 {
		return nil, errors.New("invalid length, need 256 bits")
	}
	return ToECDSA(b), nil
}

// ToECDSA creates a private key with the given D value.
func ToECDSA(prv []byte) *PrivateKey {
	if len(prv) == 0 {
		return nil
	}

	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = S256()
	priv.D = new(big.Int).SetBytes(prv)
	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(prv)
	return (*PrivateKey)(priv)
}

// Bytes returns the ecdsa PublickKey to bytes
func (pub *PublicKey) Bytes() []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(secp256k1.S256(), pub.X, pub.Y)
}

// SetBytes fill in bytes to a signature.
func (sig *Signature) SetBytes(d []byte) {
	if len(d) == 65 {
		copy(sig[:], d[:])
	}
}

// compact returns a format signature according the raw signature
func (sig *Signature) compact(data []byte, compressed bool) {
	sig.SetBytes(data)
	sig[64] += 27
	if compressed {
		sig[64] += 4
	}
}

// Ecrecover recovers publick key
func (sig *Signature) Ecrecover(hash []byte) ([]byte, error) {
	data := make([]byte, SignatureSize)
	copy(data[:], sig[:])
	data[64] = (data[64] - 27) & ^byte(4)
	return Ecrecover(hash, data[:])
}

// RecoverPublicKey recovers public key and also verifys the signature
func (sig *Signature) RecoverPublicKey(hash []byte) (*PublicKey, error) {
	s, err := sig.Ecrecover(hash)
	if err != nil {
		return nil, err
	}

	return ToECDSAPub(s), nil
}

// Verify verifys the signature with public key
func (sig *Signature) Verify(hash []byte, pub crypto.PublicKey) bool {
	sigPub, err := sig.Ecrecover(hash)
	if err != nil {
		return false
	}

	if pub != nil {
		return bytes.Equal(sigPub, pub.Bytes())
	}

	return true
}

// Bytes returns the bytes of the signature
func (sig *Signature) Bytes() []byte {
	return sig[:]
}

// VRS returns the v r s values
func (sig *Signature) VRS() (v byte, r, s *big.Int) {
	return (sig[64] - 27) & ^byte(4), new(big.Int).SetBytes(sig[:32]), new(big.Int).SetBytes(sig[32:64])
}

// Validate validates whether the signature values are valid
func (sig *Signature) Validate() bool {
	v, r, s := sig.VRS()
	one := big.NewInt(1)

	if r.Cmp(one) < 0 || s.Cmp(one) < 0 {
		return false
	}

	if s.Cmp(halfN) > 0 {
		return false
	}

	return r.Cmp(N) < 0 && s.Cmp(N) < 0 && (v == 0 || v == 1)
}

// SigToPub recovers public key from the input data to the ecdsa public key
func SigToPub(hash, sig []byte) (*PublicKey, error) {
	s, err := Ecrecover(hash, sig)
	if err != nil {
		return nil, err
	}

	return ToECDSAPub(s), nil
}

// Ecrecover recovers publick key
func Ecrecover(hash, sig []byte) ([]byte, error) {
	return secp256k1.RecoverPubkey(hash, sig)
}

// ToECDSAPub returns ecdsa public key according the input data
func ToECDSAPub(pub []byte) *PublicKey {
	if len(pub) == 0 {
		return nil
	}
	x, y := elliptic.Unmarshal(S256(), pub)
	return (*PublicKey)(&ecdsa.PublicKey{Curve: S256(), X: x, Y: y})
}

// ZeroMemory clean private key.
func (priv *PrivateKey) ZeroMemory() {
	b := priv.D.Bits()
	utils.ZeroMemory(b)
}

// func PubkeyToEAddress(p ecdsa.PublicKey) common.Address {
// 	pubBytes := (PublicKey)(&p).Bytes()
// 	return common.BytesToAddress(Keccak256(pubBytes[1:])[12:])
// }
