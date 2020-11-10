package grpcs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"math/big"
	"time"
)

type Client struct {
	node string
	GRPC *client.GrpcClient
}

func NewClient(node string) (*Client, error) {
	c := new(Client)
	c.node = node
	c.GRPC = client.NewGrpcClient(node)
	err := c.GRPC.Start()
	if err != nil {
		return nil, fmt.Errorf("grpc client start error: %v", err)
	}
	return c, nil
}

func (c *Client) SetTimeout(timeout time.Duration) error {
	if c == nil {
		return errors.New("client is nil ptr")
	}
	c.GRPC = client.NewGrpcClientWithTimeout(c.node, timeout)
	err := c.GRPC.Start()
	if err != nil {
		return fmt.Errorf("grpc start error: %v", err)
	}
	return nil
}

func (c *Client) Transfer(from, to string, amount int64) (*api.TransactionExtention, error) {
	return c.GRPC.Transfer(from, to, amount)
}

func (c *Client) GetTrxBalance(addr string) (*core.Account, error) {
	return c.GRPC.GetAccount(addr)
}
func (c *Client) GetTrc20Balance(addr, contractAddress string) (*big.Int, error) {
	return c.GRPC.TRC20ContractBalance(addr, contractAddress)
}

func (c *Client) TransferTrc20(from, to, contract string, amount *big.Int, feeLimit int64) (*api.TransactionExtention, error) {
	return c.GRPC.TRC20Send(from, to, contract, amount, feeLimit)
}

func (c *Client) BroadcastTransaction(transaction *core.Transaction) error {
	result, err := c.GRPC.Broadcast(transaction)
	if err != nil {
		return fmt.Errorf("broadcast transaction error: %v", err)
	}
	if result.Code != 0 {
		return fmt.Errorf("bad transaction: %v", string(result.GetMessage()))
	}
	if result.Result == true {
		return nil
	}
	d, _ := json.Marshal(result)
	return fmt.Errorf("tx send fail: %s", string(d))
}
