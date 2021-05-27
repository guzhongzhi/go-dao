package dao

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kisielk/sqlstruct"
	"github.com/pinguo-icc/kratos-library/logger"
	"reflect"
	"strings"
)

type SQLDAO interface {
	DAO
	DB() *sql.DB
	Table() string
	Query(sq string, params []interface{}, tx TransactionOptions) (*sql.Rows, error)
	Exec(sq string, params []interface{}, tx TransactionOptions) (sql.Result, error)
}

func NewSQLDAO(db *sql.DB, table string, idFieldName string, opts SQLDAOOptions) SQLDAO {
	if opts.Logger == nil {
		opts.Logger = logger.Default()
	}
	return &mysql{
		db:          db,
		table:       table,
		idFieldName: idFieldName,
		opts:        opts,
	}
}

type mysql struct {
	db          *sql.DB
	table       string
	idFieldName string
	opts        SQLDAOOptions
}

func (s *mysql) Table() string {
	return s.table
}

func (s *mysql) BeginTransaction(ctx context.Context, tx TxOptions) (interface{}, error) {
	var txOptions *sql.TxOptions
	if tx != nil {
		_, ok := tx.(*sql.TxOptions)
		if !ok {
			return nil, fmt.Errorf("invalid tx options for SQLDAO, must be *sql.TxOptions")
		}
		txOptions = tx.(*sql.TxOptions)
	} else {
		txOptions = &sql.TxOptions{}
	}

	return s.DB().BeginTx(context.Background(), txOptions)
}

func (s *mysql) DB() *sql.DB {
	return s.db
}

func (s *mysql) Insert(entity Data, opts InsertOptions) (interface{}, error) {
	fieldNames, params := s.buildData(entity)
	sq := "INSERT INTO `%s` (%s) VALUES (%s)"
	placeHolder := strings.Repeat("? , ", len(fieldNames))
	sq = fmt.Sprintf(sq, s.table, strings.Join(fieldNames, ", "), placeHolder[:len(placeHolder)-2])

	rs, err := s.Exec(sq, params, opts.(TransactionOptions))
	if err != nil {
		return nil, err
	}
	return rs.LastInsertId()
}

func (s *mysql) buildData(data interface{}) ([]string, []interface{}) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := reflect.TypeOf(data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fieldNames := make([]string, 0)
	params := make([]interface{}, 0)

	fieldNumber := t.NumField()
	for i := 0; i < fieldNumber; i++ {
		sqlFieldName := t.Field(i).Tag.Get("sql")
		if sqlFieldName == "" || sqlFieldName == "-" {
			continue
		}
		fieldNames = append(fieldNames, sqlFieldName)
		params = append(params, v.Field(i).Interface())
	}
	return fieldNames, params
}

func (s *mysql) Update(id interface{}, data Data, opts UpdateOptions) error {
	fieldNames, params := s.buildData(data)
	sq := "UPDATE `%s` SET %s WHERE `%s` = ?"
	setValue := strings.Join(fieldNames, " = ? ,")

	sq = fmt.Sprintf(sq, s.table, setValue[:len(setValue)-1], s.idFieldName)
	params = append(params, id)

	_, err := s.Exec(sq, params, opts.(TransactionOptions))
	if err != nil {
		return err
	}
	return nil
}

func (s *mysql) Exec(sq string, params []interface{}, tx TransactionOptions) (sql.Result, error) {
	var rs sql.Result
	var err error

	if tx != nil && tx.Tx() != nil {
		rs, err = tx.Tx().(*sql.Tx).Exec(sq, params...)
	} else {
		rs, err = s.db.Exec(sq, params...)
	}
	return rs, err
}

func (s *mysql) Find(data interface{}, opts FindOptions) error {
	filter, err := opts.Filter()

	if err != nil {
		return err
	}
	sq := ""
	if s.opts.FindSQL != "" {
		sq = s.opts.FindSQL
	}

	sq = fmt.Sprintf(sq, filter)
	rs, err := s.Query(sq, []interface{}{}, opts.(TransactionOptions))
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(data)
	if rv.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("mysql find data should be pointer of struct or map for the index of '%s'", s.table)
	}
	rv = rv.Elem()
	if !rv.CanSet() {
		return fmt.Errorf("mysql find data can not be set for index of '%s'", s.table)
	}

	if rv.Type().Kind() != reflect.Slice {
		return fmt.Errorf("the data is not a slice")
	}

	for rs.Next() {
		var item interface{}
		nVal := reflect.New(rv.Type().Elem())
		if rv.Type().Elem().Kind() == reflect.Struct {
			item = nVal.Interface()
		} else if rv.Type().Elem().Kind() == reflect.Ptr {
			item = nVal.Elem().Interface()
			t := reflect.TypeOf(item).Elem()
			item = reflect.New(t).Interface()
		}

		err := sqlstruct.Scan(item, rs)
		if err != nil {
			return err
		}
		if rv.Type().Elem().Kind() == reflect.Struct {
			rv.Set(reflect.Append(rv, reflect.ValueOf(nVal.Elem().Interface())))
		} else {
			rv.Set(reflect.Append(rv, reflect.ValueOf(item)))
		}
	}
	return nil
}

func (s *mysql) Delete(id interface{}, opts DeleteOptions) error {
	sq := "DELETE FROM `%s` WHERE `%s` = ?"
	sq = fmt.Sprintf(sq, s.table, s.idFieldName)
	_, err := s.Exec(sq, []interface{}{id}, opts.(TransactionOptions))
	if err != nil {
		return err
	}
	return nil
}

func (s *mysql) newStruct(data interface{}) interface{} {
	var ins interface{}
	t := reflect.TypeOf(data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		ins = reflect.New(t).Elem().Interface()
	}
	return ins
}

func (s *mysql) FindOne(id interface{}, data Data, opts FindOptions) error {
	sq := ""
	ins := s.newStruct(data)

	if s.opts.GetSQL != "" {
		sq = s.opts.GetSQL
		sq = fmt.Sprintf(sq, sqlstruct.Columns(ins))
	} else {
		sq = fmt.Sprintf("SELECT %s FROM `%s` WHERE `%s` = ?", sqlstruct.Columns(ins), s.table, s.idFieldName)
	}

	params := []interface{}{id}
	rs, err := s.Query(sq, params, opts.(TransactionOptions))
	if err != nil {
		return err
	}
	for rs.Next() {
		err = sqlstruct.Scan(data, rs)
		break
	}
	return err
}

func (s *mysql) Query(sq string, params []interface{}, tx TransactionOptions) (*sql.Rows, error) {
	var rs *sql.Rows
	var err error

	s.opts.Logger.Debug("tx", tx.Tx())
	if tx != nil && tx.Tx() != nil {
		s.opts.Logger.Debug("query with context", tx.Tx())
		rs, err = tx.Tx().(*sql.Tx).Query(sq, params...)
	} else {
		rs, err = s.db.Query(sq, params...)
	}
	if err != nil {
		return nil, err
	}
	return rs, err
}
