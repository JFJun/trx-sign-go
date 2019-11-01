package utils

import (
	"math/big"
)

// eth etc address
// BytesToAddress(Keccak256(pubBytes[1:])[12:])

// bitcoin litecoin bcc address
// Ripemd160Sha256(pubBytes)

const (
	// AddressSize represents the real address size.
	AddressSize = 20
)

// Address represents the 20 byte address of an Ethereum account.
// also represnets the hash160.
type Address [AddressSize]byte

// BytesToAddress converts bytes to address.
func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

func StringToAddress(s string) Address { return BytesToAddress([]byte(s)) }
func BigToAddress(b *big.Int) Address  { return BytesToAddress(b.Bytes()) }
func HexToAddress(s string) Address    { return BytesToAddress(HexToBytes(s)) }

// SetBytes sets the address to the value of b. If b is larger than len(a) it will panic
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressSize:]
	}
	copy(a[AddressSize-len(b):], b)
}
