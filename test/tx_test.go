package test

import (
	"encoding/json"
	"fmt"
	"github.com/JFJun/trx-sign-go/grpcs"
	"github.com/JFJun/trx-sign-go/sign"
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
	acc, err := c.GetTrxBalance("TFXf56UG1bhWkZq7WQEf7XW5hZXku17E8M")
	if err != nil {
		t.Fatal(err)
	}
	d, _ := json.Marshal(acc)
	fmt.Println(string(d))
	fmt.Println(acc.GetBalance())
	amount, err := c.GetTrc20Balance("TFXf56UG1bhWkZq7WQEf7XW5hZXku17E8M", "TLBaRhANQoJFTqre9Nf1mjuwNWjCJeYqUL")
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
