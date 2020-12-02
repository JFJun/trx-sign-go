package test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/JFJun/trx-sign-go/grpcs"
	"github.com/JFJun/trx-sign-go/sign"
	"github.com/btcsuite/btcutil/base58"
	addr "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"math/big"
	"testing"
)

func Test_TransferTrx(t *testing.T) {
	c, err := grpcs.NewClient("47.252.19.181:50051")
	if err != nil {
		t.Fatal(err)
	}
	tx, err := c.Transfer("TFXf56UG1bhWkZq7WQEf7XW5hZXku17E8M", "TThCjw3z2QYt4G8pgAAcVT1JhvSFBrH4U5", 1000000)
	if err != nil {
		t.Fatal(err)
	}
	signTx, err := sign.SignTransaction(tx.Transaction, "e57fa312d8e1cc656aa106e854bb7cb108ad2f77d4f99c70f3a162fdbcb682d2")
	if err != nil {
		t.Fatal(err)
	}
	err = c.BroadcastTransaction(signTx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.BytesToHexString(tx.GetTxid()))

}

func Test_GetBalance(t *testing.T) {
	c, err := grpcs.NewClient("47.252.19.181:50051")
	if err != nil {
		t.Fatal(err)
	}
	acc, err := c.GetTrxBalance("TQ54dbtZgzxbde1oWCKjQznuijLdwHubh6")
	if err != nil {
		t.Fatal(err)
	}
	d, _ := json.Marshal(acc)
	fmt.Println(string(d))
	fmt.Println(acc.GetBalance())

}

func Test_GetTrc20Balance(t *testing.T) {
	c, err := grpcs.NewClient("47.252.19.181:50051")
	if err != nil {
		t.Fatal(err)
	}
	amount, err := c.GetTrc20Balance("TK1UXQBkvAwBypz1bTWcuLHFaB8JmTjoUw", "TWVVcRqRmpyAi9dASvTXrqnS7FrwvDezMn")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(amount.String())

}

func Test_TransferTrc20(t *testing.T) {
	c, err := grpcs.NewClient("47.252.19.181:50051")
	if err != nil {
		t.Fatal(err)
	}
	tx, err := c.TransferTrc20("TFXf56UG1bhWkZq7WQEf7XW5hZXku17E8M", "TThCjw3z2QYt4G8pgAAcVT1JhvSFBrH4U5",
		"TLBaRhANQoJFTqre9Nf1mjuwNWjCJeYqUL", big.NewInt(1000000000000000000), 100000000)
	signTx, err := sign.SignTransaction(tx.Transaction, "00000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Fatal(err)
	}
	err = c.BroadcastTransaction(signTx)
	if err != nil {
		t.Fatal(err)

	}
	fmt.Println(common.BytesToHexString(tx.GetTxid()))
}

func Test_TransferTrc10(t *testing.T) {
	c, err := grpcs.NewClient("47.252.19.181:50051")
	if err != nil {
		t.Fatal(err)
	}
	from, _ := addr.Base58ToAddress("TFXf56UG1bhWkZq7WQEf7XW5hZXku17E8M")
	to, _ := addr.Base58ToAddress("TL4ebGiBbBPjduKaNEoPytVyzEuPEsFYz9")
	tokenID := "1000016"
	tx, err := c.GRPC.TransferAsset(from.String(), to.String(), tokenID, int64(123456))
	signTx, err := sign.SignTransaction(tx.Transaction, "")
	if err != nil {
		t.Fatal(err)
	}
	err = c.BroadcastTransaction(signTx)
	if err != nil {
		t.Fatal(err)

	}
	fmt.Println(common.BytesToHexString(tx.GetTxid()))
}

func Test_GetTrc10Balance(t *testing.T) {
	//c, err := grpcs.NewClient("47.252.19.181:50051")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//amount, err := c.GetTrc10Balance("TFXf56UG1bhWkZq7WQEf7XW5hZXku17E8M", "1000016")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(amount)
	_, err := DecodeCheck("TFXf56UG1bhWkZq7WQEf7XW5hZXku17E8M")
	if err != nil {
		t.Fatal(err)
	}

}

func DecodeCheck(input string) ([]byte, error) {
	decodeCheck := base58.Decode(input)
	if len(decodeCheck) == 0 {
		return nil, fmt.Errorf("b58 decode %s error", input)
	}

	if len(decodeCheck) < 4 {
		return nil, fmt.Errorf("b58 check error")
	}

	decodeData := decodeCheck[:len(decodeCheck)-4]

	h256h0 := sha256.New()
	h256h0.Write(decodeData)
	h0 := h256h0.Sum(nil)

	h256h1 := sha256.New()
	h256h1.Write(h0)
	h1 := h256h1.Sum(nil)

	if h1[0] == decodeCheck[len(decodeData)] &&
		h1[1] == decodeCheck[len(decodeData)+1] &&
		h1[2] == decodeCheck[len(decodeData)+2] &&
		h1[3] == decodeCheck[len(decodeData)+3] {
		return decodeData, nil
	}
	return nil, fmt.Errorf("b58 check error")
}

func Test_GetBlock(t *testing.T) {
	c, err := grpcs.NewClient("47.252.19.181:50051")
	if err != nil {
		t.Fatal(err)
	}
	lb, err := c.GRPC.GetNowBlock()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(lb.BlockHeader.RawData.Number)
	fmt.Println(hex.EncodeToString(lb.Blockid))
}

func Test_GetTxByTxid(t *testing.T) {
	c, err := grpcs.NewClient("47.252.19.181:50051")
	if err != nil {
		t.Fatal(err)
	}
	ti, err := c.GRPC.GetTransactionInfoByID("d81e91611935b67dd7754e107fc73c76a90b6bb20899f3450e6bf5f7b3804ac3")
	if err != nil {
		t.Fatal(err)
	}
	fee := ti.Receipt.GetEnergyFee() + ti.Receipt.GetNetFee()
	fmt.Println(fee)
}
