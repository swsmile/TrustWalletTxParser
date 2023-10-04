package main

import (
	"awesomeProject/parser"
	"awesomeProject/repository/mem_store"
	"awesomeProject/repository/rpc"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	store := mem_store.NewInMemoryStore()
	myRPCClient := rpc.NewRpcClient()
	p := parser.NewETHParser(store, myRPCClient)

	// TODO I didn't have time to implement the following part in a clean-architecture way
	in := bufio.NewScanner(os.Stdin)
	for in.Scan() {
		switch in.Text() {
		case "exit":
			return
		case "block":
			block := p.GetCurrentBlock()
			log.Println(fmt.Sprintf("Last processed block: %d", block))
		case "tx":
			log.Println("Enter address:")
			if in.Scan() {
				address := in.Text()
				log.Println(fmt.Sprintf("Transactions for address %s:", address))
				transactions := p.GetTransactions(address)
				if len(transactions) == 0 {
					log.Println("Not found")
				} else {
					marshal, _ := json.Marshal(transactions)
					log.Println(fmt.Sprintf("%s", string(marshal)))
				}
			}
		case "sub":
			log.Println("Enter address:")
			if in.Scan() {
				address := in.Text()
				if p.Subscribe(address) {
					log.Println(fmt.Sprintf("Address %s has been subscribed", address))
				} else {
					log.Println(fmt.Sprintf("Address %s has already been subscribed", address))
				}
			}
		default:
			log.Println("Unknown command")
		}
	}

	// TODO listen system signal to gracefully shutdown. Although in the current case, it's not necessary
}
