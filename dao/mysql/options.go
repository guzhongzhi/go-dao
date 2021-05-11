package mysql

import (
	"database/sql"
	daoOptions "github.com/guzhongzhi/gmicro/dao/options"
)

type DAOOptions struct {
	GetSQL    string
	FindSQL   string
	DeleteSQL string
	UpdateSQL string
	InsertSQL string
}

type TransOptions interface {
	SetTx(t *sql.Tx)
	Tx() *sql.Tx
}

type TransOptionsDefault struct {
	daoOptions.FindOptions
	tx *sql.Tx
}

func (s *TransOptionsDefault) SetTx(t *sql.Tx) {
	s.tx = t
}

func (s *TransOptionsDefault) Tx() *sql.Tx {
	return s.tx
}

func (s *TransOptionsDefault) Options() (interface{}, error) {
	return s, nil
}

type FindOptions struct {
	TransOptionsDefault
}

type InsertOptions struct {
	TransOptionsDefault
}

type UpdateOptions struct {
	TransOptionsDefault
}

type DeleteOptions struct {
	TransOptionsDefault
}

type GetOptions struct {
	TransOptionsDefault
}
