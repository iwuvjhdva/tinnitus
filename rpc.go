package tinnitus

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
)

type RPCClient struct {
	*rpcclient.Client
}

var RPC *RPCClient

func InitRPC() {
	connCfg := &rpcclient.ConnConfig{
		Host:         Config.RPC.Host,
		User:         Config.RPC.User,
		Pass:         Config.RPC.Password,
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	client, err := rpcclient.New(connCfg, nil)

	RPC = &RPCClient{client}

	if err != nil {
		Logger.Fatal("Error creating RPC client: ", err)
	}
}

func (client *RPCClient) GetBitcoinBlockVerboseTx(blockHash *chainhash.Hash) (*GetBitcoinBlockVerboseResult, error) {
	var err error
	var rawMessage json.RawMessage
	var result GetBitcoinBlockVerboseResult

	params := []json.RawMessage{
		json.RawMessage(fmt.Sprintf(`"%s"`, blockHash)),
		json.RawMessage(`2`),
	}

	rawMessage, err = client.RawRequest("getblock", params)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rawMessage, &result)

	return &result, nil
}

func (client *RPCClient) GetZcashBlockVerboseTx(blockHeight int64) (*GetZcashBlockVerboseResult, error) {
	var err error
	var rawMessage json.RawMessage
	var result GetZcashBlockVerboseResult

	params := []json.RawMessage{
		json.RawMessage(fmt.Sprintf(`"%d"`, blockHeight)),
		json.RawMessage(`2`),
	}

	rawMessage, err = client.RawRequest("getblock", params)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rawMessage, &result)

	return &result, nil
}

type GetBitcoinBlockVerboseResult struct {
	*btcjson.GetBlockVerboseResult
	Tx []btcjson.TxRawResult `json:"tx,omitempty"`
}

type GetZcashBlockVerboseResult struct {
	Hash          string                `json:"hash"`
	Confirmations int64                 `json:"confirmations"`
	StrippedSize  int32                 `json:"strippedsize"`
	Size          int32                 `json:"size"`
	Weight        int32                 `json:"weight"`
	Height        int64                 `json:"height"`
	Version       int32                 `json:"version"`
	VersionHex    string                `json:"versionHex"`
	MerkleRoot    string                `json:"merkleroot"`
	Tx            []btcjson.TxRawResult `json:"tx,omitempty"`
	Time          int64                 `json:"time"`
	Nonce         string                `json:"nonce"`
	Bits          string                `json:"bits"`
	Difficulty    float64               `json:"difficulty"`
	PreviousHash  string                `json:"previousblockhash"`
	NextHash      string                `json:"nextblockhash,omitempty"`
}
