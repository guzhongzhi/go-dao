package dao

type GetOptions interface {
	Options() (interface{}, error)
}

type NoopTransactionOptions struct {
}

func (s *NoopTransactionOptions) SetTx(v interface{}) {

}

func (s *NoopTransactionOptions) Tx() interface{} {
	return nil
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
