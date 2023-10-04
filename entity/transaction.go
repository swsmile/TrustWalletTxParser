// Package entity defines all entities, which don't couple with a specified storage model
package entity

// Transaction is defined based on https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_gettransactionbyhash
// Read through https://ethereum.org/en/developers/docs/intro-to-ethereum/#transactions for more details of Transaction
type Transaction struct {
	// Note: https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_gettransactionbyhash doesn't mention the valid range of integer fields and I just assume int64 is good enough
	BlockHash        string
	BlockNumber      int64
	From             string
	Gas              int64
	GasPrice         int64
	Hash             string
	Input            string
	Nonce            int64
	To               string
	TransactionIndex int64

	// The following fields are supposed to be an integer as per https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_gettransactionbyhash, but they may be larger than max_int64
	Value string
	V     string
	R     string
	S     string
}
