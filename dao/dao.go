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
	Find(data interface{}, opts FindOptions) error
	Delete(id interface{}, opts DeleteOptions) error
	Get(id interface{}, data Entity, opts GetOptions) error
}
