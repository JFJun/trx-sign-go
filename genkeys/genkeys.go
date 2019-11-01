package genkeys

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec"
	"github.com/jfjun/wallet-tools/crypto"
	"github.com/jfjun/wallet-tools/encoding/base58"
	"math/big"
)

/*
trx 生成私钥和地址
*/




/*
随机生成私钥，返回32字节私钥和64字节的公钥
*/
func GenerateTrxKeys()([]byte,[]byte){
	privkey,err:=btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		panic(err)
	}
	priv := privkey.D.Bytes()
	//避免priv的len不是32
	if len(priv)!=32{
		for true{
			newPrivKey,err:=btcec.NewPrivateKey(btcec.S256())
			if err != nil {
				//if have some error ,cut this exe
				panic(err)
			}
			priv = newPrivKey.D.Bytes()
			if len(priv)==32{
				break
			}
		}
	}
	//未压缩公钥
	pub:=privkey.PubKey().SerializeUncompressed()
	//和btc的唯一差别是不需要加上 04
	return priv,pub[1:]
}


/*
根据私钥生成公钥
*/
func PrivateBytesToPub(prv []byte)(pub []byte){
	k:=new(big.Int).SetBytes(prv)
	c:=btcec.S256()
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve =c
	priv.D = k
	priv.PublicKey.X, priv.PublicKey.Y = c.ScalarBaseMult(k.Bytes())
	btcPriv:=(*btcec.PrivateKey)(priv)
	pub =btcPriv.PubKey().SerializeUncompressed()
	return pub[1:]
}

/*
64字节公钥转化为地址
*/
func PubToAddress(pub []byte)(b58Addr,hexAddr string){
	data:=crypto.Keccak256(pub)
	d:=data[11:]
	d[0] = 0x41
	hexAddr = hex.EncodeToString(d)
	checkSum:=crypto.DoubleSha256(d)
	d = append(d, checkSum[:4]...)
	b58Addr = base58.Encode(d)
	return b58Addr,hexAddr
}

func TrxHexAddressToBase58(hexAddress string)string{
	data,err:=hex.DecodeString(hexAddress)
	if err != nil {
		panic(err)
	}
	checkSum:=crypto.DoubleSha256(data)
	data = append(data,checkSum[:4]...)
	return base58.Encode(data)
}

func TrxBase58AddressToHex(b58Address string)string{
	data:=base58.Decode(b58Address)
	return hex.EncodeToString(data[:len(data)-4])
}