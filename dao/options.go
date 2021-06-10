package dao

type TxOptions interface {
}

type NoopTransactionOptions struct {
	tx interface{}
}

func (s *NoopTransactionOptions) SetTx(v interface{}) {
	s.tx = v
}

func (s *NoopTransactionOptions) Tx() interface{} {
	return s.tx
}

type TransactionOptions interface {
	SetTx(tx interface{})
	Tx() interface{}
}

type FindOptions interface {
	TransactionOptions
	Filter() (interface{}, error)
	Options() (interface{}, error)
	Pagination() *Pagination
	Sorts() interface{}
}

type InsertOptions interface {
	TransactionOptions
	Options() (interface{}, error)
}

type UpdateOptions interface {
	TransactionOptions
	Options() (interface{}, error)
}

type DeleteOptions interface {
	TransactionOptions
	Options() (interface{}, error)
}

type CollectionOptions interface {
	Options() (interface{}, error)
}
