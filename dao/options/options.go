package options

import "github.com/guzhongzhi/gmicro/dao/pagination"

type GetOptions interface {
	Options() (interface{}, error)
}

type FindOptions interface {
	Filter() (interface{}, error)
	Options() (interface{}, error)
	Pagination() *pagination.Pagination
}

type InsertOptions interface {
	Options() (interface{}, error)
}

type UpdateOptions interface {
	Options() (interface{}, error)
}

type DeleteOptions interface {
	Options() (interface{}, error)
}

type CollectionOptions interface {
	Options() (interface{}, error)
}
