package dao

type Entity interface {
	IsNew() bool
	ID() interface{}
	SetID(v interface{})
	String() string
}

type DAO interface {
	Insert(entity Entity, opts InsertOptions) (id interface{}, err error)
	Update(id interface{}, data Entity, opts UpdateOptions) error
	Find(opts interface{}, data interface{}) error
	Delete(id interface{}) error
	Get(id interface{}, data Entity) error
}

type FindOptions interface {
	Filter() (interface{}, error)
	Options() (interface{}, error)
	Pagination() *Pagination
}

type InsertOptions interface {
	Options() interface{}
}

type UpdateOptions interface {
	Options() interface{}
}

type CollectionOptions interface {
	Options() interface{}
}
