package parser

import (
	"awesomeProject/util/log"
	"sync"
	"time"
)

const (
	initMostRecentBlockNumber = -1 // -1 means no blocks parsed yet. I read through https://ethereum.org/en/developers/docs/apis/json-rpc/, but didn't figure out the valid range of block number. So I just use -1 to indicate no blocks parsed yet with the assumption that valid block numbers are always positive.

	daemonTickIntervalInSecond = 1 * time.Second
)

type daemon struct {
	mostRecentBlockNumber        int64
	mostRecentFetchedBlockNumber int64

	subscribersSet map[string]interface{}
	lock           sync.RWMutex

	store  Store
	client RpcClient
}

func newDaemon(store Store, rpcClient RpcClient) *daemon {
	return &daemon{
		mostRecentBlockNumber:        initMostRecentBlockNumber,
		mostRecentFetchedBlockNumber: initMostRecentBlockNumber,

		subscribersSet: make(map[string]interface{}),
		lock:           sync.RWMutex{},

		store:  store,
		client: rpcClient,
	}
}

func (d *daemon) run() {
	for {
		d.tick()
		time.Sleep(daemonTickIntervalInSecond)
	}
}

func (d *daemon) tick() {
	newMostRecentBlockNumber, err := d.client.GetMostRecentBlockNumber()
	if err != nil {
		log.Warning("error getting recent block number, err: %v", err)
		return
	}

	d.lock.Lock()
	defer d.lock.Unlock()

	// TODO sometimes, a new block number is returned, but the block is not available yet. We could optimize the implementation a bit.
	if d.mostRecentBlockNumber == initMostRecentBlockNumber {
		err = d.fetchBlockTransactionsByBlockNum(newMostRecentBlockNumber)
		if err != nil {
			return
		}
	} else {
		for blockNumber := d.mostRecentFetchedBlockNumber + 1; blockNumber <= newMostRecentBlockNumber; blockNumber++ {
			err = d.fetchBlockTransactionsByBlockNum(blockNumber)
			if err != nil {
				return
			}
		}
	}

	d.mostRecentBlockNumber = newMostRecentBlockNumber
	d.mostRecentFetchedBlockNumber = newMostRecentBlockNumber
}

func (d *daemon) fetchBlockTransactionsByBlockNum(blockNumber int64) error {
	transactions, err := d.client.GetBlockTransactionsByBlockNumber(blockNumber)
	if err != nil {
		log.Warning("error on getting blockNumber by number, err: %v", err)
		return err
	}

	for _, transaction := range transactions {
		// Inbound and outbound transactions are not needed to differentiate as per the requirement
		if _, exist := d.subscribersSet[transaction.To]; exist {
			d.store.Insert(transaction.To, transaction)
		}
		if _, exist := d.subscribersSet[transaction.From]; exist {
			d.store.Insert(transaction.From, transaction)
		}
	}

	return nil
}

func (d *daemon) mostRecentParsedBlockNumber() int64 {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.mostRecentBlockNumber
}

func (d *daemon) subscribe(address string) bool {
	d.lock.Lock()
	defer d.lock.Unlock()
	if _, exist := d.subscribersSet[address]; exist {
		return false
	}

	d.subscribersSet[address] = struct{}{}
	return true
}
