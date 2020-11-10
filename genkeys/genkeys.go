package genkeys

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	addr "github.com/fbsobreira/gotron-sdk/pkg/address"
)

func GenerateKey() (wif string, address string) {
	priv, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return "", ""
	}
	if len(priv.D.Bytes()) != 32 {
		for {
			priv, err := btcec.NewPrivateKey(btcec.S256())
			if err != nil {
				continue
			}
			if len(priv.D.Bytes()) == 32 {
				break
			}
		}
	}
	a := addr.PubkeyToAddress(priv.ToECDSA().PublicKey)
	address = a.String()
	wif = hex.EncodeToString(priv.D.Bytes())
	return
}

func CreateAddressBySeed(seed []byte) (string, error) {
	if len(seed) != 32 {
		return "", fmt.Errorf("seed len=[%d] is not equal 32", len(seed))
	}
	priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), seed)
	if priv == nil {
		return "", errors.New("priv is nil ptr")
	}
	a := addr.PubkeyToAddress(priv.ToECDSA().PublicKey)
	return a.String(), nil
}

func AddressB58ToHex(b58 string) (string, error) {
	a, err := addr.Base58ToAddress(b58)
	if err != nil {
		return "", err
	}
	return a.Hex(), nil
}

func AddressHexToB58(hexAddress string) string {
	a := addr.HexToAddress(hexAddress)
	return a.String()
}
