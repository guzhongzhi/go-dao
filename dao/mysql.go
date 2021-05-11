package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	mysql2 "github.com/guzhongzhi/gmicro/dao/mysql"
	"github.com/guzhongzhi/gmicro/dao/options"
	"github.com/kisielk/sqlstruct"
	"reflect"
	"strings"
)

type MysqlDAO interface {
	DAO
	DB() *sql.DB
	Table() string
	Query(sq string, params []interface{}, tx *sql.Tx) (*sql.Rows, error)
	Exec(sq string, params []interface{}, tx *sql.Tx) (sql.Result, error)
}

func NewMysqlDAO(db *sql.DB, table string, idFieldName string, opts mysql2.DAOOptions) MysqlDAO {
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
	opts        mysql2.DAOOptions
}

func (s *mysql) Table() string {
	return s.table
}

func (s *mysql) DB() *sql.DB {
	return s.db
}

func (s *mysql) Insert(entity Entity, opts options.InsertOptions) (interface{}, error) {
	if opts == nil {
		opts = &mysql2.InsertOptions{}
	}

	fieldNames, params := s.buildData(entity)
	sq := "INSERT INTO `%s` (%s) VALUES (%s)"
	placeHolder := strings.Repeat("? , ", len(fieldNames))
	sq = fmt.Sprintf(sq, s.table, strings.Join(fieldNames, ", "), placeHolder[:len(placeHolder)-2])

	rs, err := s.Exec(sq, params, opts.(mysql2.TransOptions).Tx())
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

func (s *mysql) Update(id interface{}, data Entity, opts options.UpdateOptions) error {
	if opts != nil {
		opts = &mysql2.UpdateOptions{}
	}

	fieldNames, params := s.buildData(data)
	sq := "UPDATE `%s` SET %s WHERE `%s` = ?"
	setValue := strings.Join(fieldNames, " = ? ,")

	sq = fmt.Sprintf(sq, s.table, setValue[:len(setValue)-1], s.idFieldName)
	params = append(params, id)

	_, err := s.Exec(sq, params, opts.(mysql2.TransOptions).Tx())
	if err != nil {
		return err
	}
	return nil
}

func (s *mysql) Exec(sq string, params []interface{}, tx *sql.Tx) (sql.Result, error) {
	var rs sql.Result
	var err error

	if tx != nil {
		rs, err = tx.Exec(sq, params...)
	} else {
		rs, err = s.db.Exec(sq, params...)
	}
	return rs, err
}

func (s *mysql) Find(data interface{}, opts options.FindOptions) error {
	if opts == nil {
		opts = &mysql2.FindOptions{}
	}
	fmt.Println("opts.(mysql2.TransOptions).Tx()", opts.(mysql2.TransOptions).Tx())
	filter, err := opts.Filter()

	if err != nil {
		return err
	}
	sq := ""
	if s.opts.FindSQL != "" {
		sq = s.opts.FindSQL
	}

	sq = fmt.Sprintf(sq, filter)
	rs, err := s.Query(sq, []interface{}{}, opts.(mysql2.TransOptions).Tx())
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

func (s *mysql) Delete(id interface{}, opts options.DeleteOptions) error {
	if opts == nil {
		opts = &mysql2.DeleteOptions{}
	}
	o := opts.(mysql2.TransOptions)

	sq := "DELETE FROM `%s` WHERE `%s` = ?"
	sq = fmt.Sprintf(sq, s.table, s.idFieldName)
	_, err := s.Exec(sq, []interface{}{id}, o.Tx())
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

func (s *mysql) Get(id interface{}, data Entity, opts options.GetOptions) error {
	if opts == nil {
		opts = &mysql2.GetOptions{}
	}
	o := opts.(mysql2.TransOptions)

	sq := ""
	ins := s.newStruct(data)

	if s.opts.GetSQL != "" {
		sq = s.opts.GetSQL
		sq = fmt.Sprintf(sq, sqlstruct.Columns(ins))
	} else {
		sq = fmt.Sprintf("SELECT %s FROM `%s` WHERE `%s` = ?", sqlstruct.Columns(ins), s.table, s.idFieldName)
	}

	params := []interface{}{id}
	rs, err := s.Query(sq, params, o.Tx())
	if err != nil {
		return err
	}
	for rs.Next() {
		err = sqlstruct.Scan(data, rs)
		break
	}
	return err
}

func (s *mysql) Query(sq string, params []interface{}, tx *sql.Tx) (*sql.Rows, error) {
	var rs *sql.Rows
	var err error

	if tx != nil {
		rs, err = tx.Query(sq, params...)
	} else {
		rs, err = s.db.Query(sq, params...)
	}
	if err != nil {
		return nil, err
	}
	return rs, err
}
