package dao

import "context"

type Entity interface {
	GetData() Data
	SetData(v Data)
}

type Data interface {
	IsNew() bool
	GetID() interface{}
	SetID(v interface{})
	String() string
}

type DAO interface {
	Insert(entity Data, opts InsertOptions) (id interface{}, err error)
	Update(id interface{}, data Data, opts UpdateOptions) error
	Find(data interface{}, opts FindOptions) error
	Delete(id interface{}, opts DeleteOptions) error
	FindOne(id interface{}, data Data, opts FindOptions) error
	BeginTransaction(ctx context.Context, tx TxOptions) (interface{}, error)
}
