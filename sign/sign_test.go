package sign

import (
	"fmt"
	"testing"
)

func TestTrxSignByPrivateStr(t *testing.T) {
	sig,err:=TrxSignByPrivateStr("0a028f8b22088b73d8d9386a59d340c0bcddb0e22d5a69080112650a2d747970652e67" +
		"6f6f676c65617069732e636f6d2f70726f746f636f6c2e5472616e73666572436f6e747261637412340a154106c6a577814bc4f1cf7305" +
		"7ada02a3a8e05f51481215417eb81b56a44ed4c13f3afbdc4c068bbf1b2b5d39189adbd6ee0170deeed9b0e22d","245" +
		"e817e1d8cc348b40b366b22fb6070feb6dfd2ee441d13c193114f4987d45f")
	if err != nil {
		panic(err)
	}
	fmt.Println(sig)
	/*
	print: 9ab2871cbeb65e9b96d6e8f2ebb185219fe7170b2d489f624ef32d6be657f9684ec9a4224588315605e250e4d2513dbc8fe9adc9c86deeec03a25c8bbbe4aca51c
	*/
}
