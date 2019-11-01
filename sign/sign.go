package sign

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/goleveldb/leveldb/errors"
	"github.com/jfjun/wallet-tools/crypto"
	"github.com/jfjun/wallet-tools/crypto/secp256k1"
)

/*
trx 离线签名
*/


/*
	raw_data_hex: 通过wallet/createtransaction接口生成，详细http接口文档： https://github.com/tronprotocol/documentation/blob/master/TRX_CN/Tron-http.md
	privateKeyStr：16进制的私钥
*/
func TrxSignByPrivateStr(raw_data_hex ,privateKeyStr string)(string,error){
	private_key_bytes,err:=hex.DecodeString(privateKeyStr)
	if err != nil {
		return "",fmt.Errorf("Hex decode private key error,Err=[%v]",err)
	}
	return TrxSignByPrivateBytes(raw_data_hex,private_key_bytes)
}

/*
	privateKeyBytes: 32字节私钥
*/
func TrxSignByPrivateBytes(raw_data_hex string,privateKeyBytes []byte)(string,error){
	if len(privateKeyBytes)!=32 {
		return "",errors.New("private key bytes length is not equal 32")
	}
	priv:=secp256k1.ToECDSA(privateKeyBytes)
	raw_data,err:=hex.DecodeString(raw_data_hex)
	if err != nil {
		return "",fmt.Errorf("Hex decode raw data hex error,Err=[%v]",err)
	}
	hash:=crypto.Sha256(raw_data)
	sig,err:=priv.Sign(hash[:])
	if err != nil {
		return "",fmt.Errorf(" Trx sign hash error,Err=[%v]",err)
	}
	return hex.EncodeToString(sig.Bytes()),nil
}