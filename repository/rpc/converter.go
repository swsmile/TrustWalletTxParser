package rpc

import (
	"awesomeProject/entity"
	"awesomeProject/util"
	"awesomeProject/util/log"
	"strings"
)

func parseBlockNumberFromResp(res *blockNumberResponse) (int64, error) {
	parsedBlockNumber, errP := util.HexToInt(strings.Replace(res.Result, "0x", "", -1))
	if errP != nil {
		return 0, errP
	}
	return parsedBlockNumber, nil
}

func convertDTOToTransactionEntity(transactionDTO transaction, err error) (entity.Transaction, error) {
	t := entity.Transaction{
		BlockHash: transactionDTO.BlockHash,
		From:      transactionDTO.From,
		Hash:      transactionDTO.Hash,
		Input:     transactionDTO.Input,
		To:        transactionDTO.To,

		Value: transactionDTO.Value,
		V:     transactionDTO.V,
		R:     transactionDTO.R,
		S:     transactionDTO.S,
	}

	var parsedBlockNumber, parsedGas, parsedGasPrice, parsedNonce, parsedTransactionIndex int64
	parsedBlockNumber, err = util.HexToInt(transactionDTO.BlockNumber)
	if err != nil {
		log.Warning("error on parsing block number, err: %v, transactionDTO.BlockNumber: %v", err, transactionDTO.BlockNumber)
		return entity.Transaction{}, err
	}
	t.BlockNumber = parsedBlockNumber

	parsedGas, err = util.HexToInt(transactionDTO.Gas)
	if err != nil {
		log.Warning("error on parsing gas, err: %v, transactionDTO.Gas: %v", err, transactionDTO.Gas)
		return entity.Transaction{}, err
	}
	t.Gas = parsedGas

	parsedGasPrice, err = util.HexToInt(transactionDTO.GasPrice)
	if err != nil {
		log.Warning("error on parsing gas price, err: %v, transactionDTO.GasPrice: %v", err, transactionDTO.GasPrice)
		return entity.Transaction{}, err
	}
	t.GasPrice = parsedGasPrice

	parsedNonce, err = util.HexToInt(transactionDTO.Nonce)
	if err != nil {
		log.Warning("error on parsing nonce, err: %v, transactionDTO.Nonce: %v", err, transactionDTO.Nonce)
		return entity.Transaction{}, err
	}
	t.Nonce = parsedNonce

	parsedTransactionIndex, err = util.HexToInt(transactionDTO.TransactionIndex)
	if err != nil {
		log.Warning("error on parsing t index, err: %v, transactionDTO.TransactionIndex: %v", err, transactionDTO.TransactionIndex)
		return entity.Transaction{}, err
	}
	t.TransactionIndex = parsedTransactionIndex

	//parsedValue, err = util.HexToInt(transactionDTO.Value)
	//if err != nil {
	//	log.Warning("error parsing value, err: %v, transactionDTO.Value: %v", err, transactionDTO.Value)
	//	return entity.Transaction{}, err
	//}
	//t.Value = parsedValue
	//
	//parsedV, err = util.HexToInt(transactionDTO.V)
	//if err != nil {
	//	log.Warning("error parsing v, err: %v, transactionDTO.V: %v", err, transactionDTO.V)
	//	return entity.Transaction{}, err
	//}
	//t.V = parsedV

	//parsedR, err = util.HexToInt(transactionDTO.R)
	//if err != nil {
	//	log.Warning("error parsing r, err: %v, transactionDTO.R: %v", err, transactionDTO.R)
	//	return entity.Transaction{}, err
	//}
	//t.R = parsedR

	//parsedS, err = util.HexToInt(transactionDTO.S)
	//if err != nil {
	//	log.Warning("error parsing s, err: %v, transactionDTO.S: %v", err, transactionDTO.S)
	//	return entity.Transaction{}, err
	//}
	//t.S = parsedS

	return t, nil
}
