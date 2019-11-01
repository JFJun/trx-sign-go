package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"github.com/jfjun/wallet-tools/utils"

	"golang.org/x/crypto/scrypt"
)

// a part of ethereum.

const (
	scryptR     = 6
	scryptN     = 1 << 2
	scryptP     = 1
	scryptDKLen = 32
)

var (
	ErrDecrypt = errors.New("could not decrypt key with given passphrase")
)

func aesCTRXOR(key, inText, iv []byte) ([]byte, error) {
	// AES-128 is selected due to size of encryptKey.
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(aesBlock, iv)
	outText := make([]byte, len(inText))
	stream.XORKeyStream(outText, inText)
	return outText, err
}

// GenEncryptInfo generate salt and iv for aes encryption.
func GenEncryptInfo() (salt, iv []byte) {
	salt = utils.GetRandomBytes(32)
	iv = utils.GetRandomBytes(aes.BlockSize)
	return
}

// GetDerivedKey gets aes params according the auth.
func GetDerivedKey(auth string, salt []byte) (encryptKey []byte) {
	authArray := []byte(auth)
	derivedKey, err := scrypt.Key(authArray, salt, scryptN, scryptR, scryptP, scryptDKLen)
	if err != nil {
		return nil
	}
	return derivedKey
}

// Encrypt encrypts a key using the specified scrypt parameters into bytes
// that can be decrypted later on.
func Encrypt(derivedKey []byte, key []byte, iv []byte) ([]byte, error) {
	cipherText, err := aesCTRXOR(derivedKey[:16], key, iv)
	if err != nil {
		return nil, err
	}
	return cipherText, nil
}

// Decrypt decrypts a key from bytes, returning the private key bytes.
func Decrypt(derivedKey []byte, cipherText []byte, iv, mac []byte) ([]byte, error) {
	// check mac
	if !bytes.Equal(Keccak256(derivedKey[16:32]), mac) {
		return nil, ErrDecrypt
	}

	plainText, err := aesCTRXOR(derivedKey[:16], cipherText, iv)
	return plainText, err
}

// GetDerivedKey -> encrypt all, salt, iv and mac is constant.
// GetDerivedKey -> decrypt
