package parser

import (
	"awesomeProject/entity"
)

type RpcClient interface {
	GetMostRecentBlockNumber() (int64, error)
	GetBlockTransactionsByBlockNumber(blockNumber int64) ([]entity.Transaction, error)
}
