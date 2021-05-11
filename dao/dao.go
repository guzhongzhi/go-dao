package dao

import "github.com/guzhongzhi/gmicro/dao/options"

type Entity interface {
	IsNew() bool
	ID() interface{}
	SetID(v interface{})
	String() string
}

type DAO interface {
	Insert(entity Entity, opts options.InsertOptions) (id interface{}, err error)
	Update(id interface{}, data Entity, opts options.UpdateOptions) error
	Find(data interface{}, opts options.FindOptions) error
	Delete(id interface{}, opts options.DeleteOptions) error
	Get(id interface{}, data Entity, opts options.GetOptions) error
}
