package genkeys

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestGenerateTrxKeys(t *testing.T) {
	priv,pub:=GenerateTrxKeys()
	fmt.Println("Private: ",hex.EncodeToString(priv))
	fmt.Println("Pub: ",hex.EncodeToString(pub))
}

func TestPrivateBytesToPub(t *testing.T) {
	priv,pub:=GenerateTrxKeys()
	fmt.Println("Private1: ",hex.EncodeToString(priv))
	fmt.Println("Pub1: ",hex.EncodeToString(pub))
	pub2:=PrivateBytesToPub(priv)
	fmt.Println("Pub2: ",hex.EncodeToString(pub2))
}

func TestPubToAddress(t *testing.T) {
	priv,pub:=GenerateTrxKeys()
	fmt.Println("Private: ",hex.EncodeToString(priv))
	fmt.Println("Pub: ",hex.EncodeToString(pub))
	b58Addr,hexAddr:=PubToAddress(pub)
	fmt.Println("B58Addr: ",b58Addr)
	fmt.Println("HexAddr: ",hexAddr)
}
func TestTrxBase58AddressToHex(t *testing.T) {
	priv,pub:=GenerateTrxKeys()
	fmt.Println("Private: ",hex.EncodeToString(priv))
	fmt.Println("Pub: ",hex.EncodeToString(pub))
	b58Addr,hexAddr:=PubToAddress(pub)
	fmt.Println("B58Addr1: ",b58Addr)
	fmt.Println("HexAddr1: ",hexAddr)
	hexAddr2:=TrxBase58AddressToHex(b58Addr)
	fmt.Println("HexAddr2: ",hexAddr2)
}

func TestTrxHexAddressToBase58(t *testing.T) {
	priv,pub:=GenerateTrxKeys()
	fmt.Println("Private: ",hex.EncodeToString(priv))
	fmt.Println("Pub: ",hex.EncodeToString(pub))
	b58Addr,hexAddr:=PubToAddress(pub)
	fmt.Println("B58Addr1: ",b58Addr)
	fmt.Println("HexAddr1: ",hexAddr)
	b58Addr2:=TrxHexAddressToBase58(hexAddr)
	fmt.Println("B58Addr2: ",b58Addr2)
}