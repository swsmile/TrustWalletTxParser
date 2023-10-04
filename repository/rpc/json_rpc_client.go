package rpc

import (
	"awesomeProject/entity"
	"awesomeProject/parser"
	"awesomeProject/util"
	"awesomeProject/util/log"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

const (
	normalErrCode = 0
	ethURL        = "https://cloudflare-eth.com"
	jsonRPCVer    = "2.0"

	methodToGetMostRecentBlockNumber = "eth_blockNumber"
	methodToGetBlockByBlockNumber    = "eth_getBlockByNumber"

	// TODO to add more error codes
	errCodeNotFoundBlock = -32001
)

var (
	errNotFoundBlock = errors.New("not found block number")
	errUnknown       = errors.New("unknown")
)

// jsonRPCClient is specifically for Ethereum's JSON-RPC APIs
type jsonRPCClient struct {
	url string
	seq int
}

func NewRpcClient() parser.RpcClient {
	return &jsonRPCClient{
		url: ethURL,
		seq: 0,
	}
}

func (j *jsonRPCClient) request(method string, params []interface{}) (*http.Response, error) {
	defer func() { j.seq++ }()
	req := rpcRequest{jsonRPCVer, method, params, j.seq}
	marshal, err := json.Marshal(req)
	if err != nil {
		log.Warning("error on marshaling rpc request, req: %v", req)
		return nil, err
	}

	resp, err := http.Post(j.url, "application/json", strings.NewReader(string(marshal)))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func closeResp(resp *http.Response) {
	err := resp.Body.Close()
	if err != nil {
		log.Warning("error on closing response body, err: %v", err)
	}
}

// GetMostRecentBlockNumber returns the number of the most recent block, and it is implemented based on https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_blocknumber
func (j *jsonRPCClient) GetMostRecentBlockNumber() (int64, error) {
	resp, err := j.request(methodToGetMostRecentBlockNumber, []interface{}{})
	defer closeResp(resp)
	if err != nil {
		log.Warning("error on getting recent block number, err: %v", err)
		return 0, err
	}

	var res = new(blockNumberResponse)
	if err = json.NewDecoder(resp.Body).Decode(res); err != nil {
		log.Warning("error on decoding recent block number response, err: %v", err)
		return 0, err
	}
	// TODO log with request and response
	if res.Error.Code != normalErrCode { // E.g., {"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request"}}
		log.Warning("error on getting recent block number, errMsg: %v, errCode: %v", res.Error.Message, res.Error.Code)
		return 0, errUnknown
	}

	parsedMostRecentBlockNumber, errP := parseBlockNumberFromResp(res)
	if errP != nil {
		log.Warning("error on parsing recent block number, err: %v, res.Result: %v", errP, res.Result)
		return 0, errP
	}
	log.Info("successfully got the new MostRecentBlockNumber: %v", parsedMostRecentBlockNumber)
	return parsedMostRecentBlockNumber, nil
}

// GetBlockTransactionsByBlockNumber returns the information of a block by a block number, and it is implemented based on https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_getblockbynumber
func (j *jsonRPCClient) GetBlockTransactionsByBlockNumber(blockNumber int64) ([]entity.Transaction, error) {
	resp, err := j.request(methodToGetBlockByBlockNumber, []interface{}{
		util.IntToHex(blockNumber), true,
	})
	defer closeResp(resp)
	if err != nil {
		return nil, err
	}

	var res = new(blockResponse)
	if err = json.NewDecoder(resp.Body).Decode(res); err != nil {
		return nil, err
	}
	// TODO log with request and response
	if res.Error.Code != normalErrCode { // E.g., {"jsonrpc":"2.0","error":{"code":-32700,"message":"Parse error"},"id":null}
		log.Warning("error on getting block by block number, blockNumber: %v, errMsg: %v, errCode: %v", blockNumber, res.Error.Message, res.Error.Code)
		if res.Error.Code == errCodeNotFoundBlock {
			return nil, errNotFoundBlock
		}
		return nil, errUnknown
	}

	log.Info("successfully got block by block number, blockNumber: %v", blockNumber)
	transactions := make([]entity.Transaction, 0, len(res.Result.Transactions))
	for _, transactionDTO := range res.Result.Transactions {
		t, errT := convertDTOToTransactionEntity(transactionDTO, err)
		if errT != nil { // I achieve fast-fail here, but the actual behaviour may depend on the requirement of how to handle the abnormal.

			return nil, errT
		}

		transactions = append(transactions, t)
	}
	return transactions, nil
}
