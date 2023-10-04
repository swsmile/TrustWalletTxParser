package parser

import "awesomeProject/entity"

type Store interface {
	Insert(address string, transaction entity.Transaction)
	Get(address string) []entity.Transaction
}
