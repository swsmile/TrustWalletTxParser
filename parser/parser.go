package parser

import "awesomeProject/entity"

type Parser interface {
	// GetCurrentBlock returns the latest parsed block number
	GetCurrentBlock() int64

	// Subscribe adds an address to observer and returns true if the address is successfully added and returns false if the address is already added
	Subscribe(address string) bool

	// GetTransactions returns list of inbound or outbound transactions for an address
	// For simplicity, I don't validate the format of a given address
	GetTransactions(address string) []entity.Transaction
}

type ethParser struct {
	d *daemon
}

func NewETHParser(store Store, rpcClient RpcClient) Parser {
	p := &ethParser{newDaemon(store, rpcClient)}
	go p.d.run()
	return p
}

func (p *ethParser) GetCurrentBlock() int64 {
	return p.d.mostRecentParsedBlockNumber()
}

func (p *ethParser) Subscribe(address string) bool {
	return p.d.subscribe(address)
}

func (p *ethParser) GetTransactions(address string) []entity.Transaction {
	return p.d.store.Get(address)
}
