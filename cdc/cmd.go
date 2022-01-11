package cdc

type cmdExecutor interface {
	query() error
	queryDetail() error
	resignOwner() error
	rebalanceTable() error
}

type command struct {
	arg00 string
	arg01 string
	arg02 string
	arg03 string
	arg04 string
}
