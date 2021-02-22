package grpcs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"math/big"
	"strings"
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

/*
保持连接，如果中途连接失败，就重连
*/
func (c *Client) keepConnect() error {
	_, err := c.GRPC.GetNodeInfo()
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return c.GRPC.Reconnect(c.node)
		}
		return fmt.Errorf("node connect error: %v", err)
	}
	return nil
}

func (c *Client) Transfer(from, to string, amount int64) (*api.TransactionExtention, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}
	return c.GRPC.Transfer(from, to, amount)
}

func (c *Client) GetTrc10Balance(addr, assetId string) (int64, error) {
	err := c.keepConnect()
	if err != nil {
		return 0, err
	}
	acc, err := c.GRPC.GetAccount(addr)
	if err != nil || acc == nil {
		return 0, fmt.Errorf("get %s account error: %v", addr, err)
	}
	for key, value := range acc.AssetV2 {
		if key == assetId {
			return value, nil
		}
	}
	return 0, fmt.Errorf("%s do not find this assetID=%s amount", addr, assetId)
}

func (c *Client) GetTrxBalance(addr string) (*core.Account, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}
	return c.GRPC.GetAccount(addr)
}
func (c *Client) GetTrc20Balance(addr, contractAddress string) (*big.Int, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}
	return c.GRPC.TRC20ContractBalance(addr, contractAddress)
}

func (c *Client) TransferTrc10(from, to, assetId string, amount int64) (*api.TransactionExtention, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}
	fromAddr, err := address.Base58ToAddress(from)
	if err != nil {
		return nil, fmt.Errorf("from address is not equal")
	}
	toAddr, err := address.Base58ToAddress(to)
	if err != nil {
		return nil, fmt.Errorf("to address is not equal")
	}
	return c.GRPC.TransferAsset(fromAddr.String(), toAddr.String(), assetId, amount)
}

func (c *Client) TransferTrc20(from, to, contract string, amount *big.Int, feeLimit int64) (*api.TransactionExtention, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}
	return c.GRPC.TRC20Send(from, to, contract, amount, feeLimit)
}

func (c *Client) BroadcastTransaction(transaction *core.Transaction) error {
	err := c.keepConnect()
	if err != nil {
		return err
	}
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
