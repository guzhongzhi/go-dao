package data

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
	"time"
)

func TestTemplate(t *testing.T) {

	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", "root", "AyXwZhBf0nGlKxJf", "tcp", "127.0.0.1", 3306, "bmall")
	DB, err := sql.Open("mysql", dsn)

	if err != nil {
		panic(err)
	}
	o := &TOrder{}
	torder := NewTOrderDAO(DB)

	tx, err := torder.BeginTransaction(context.Background(), nil)
	//sqlTx := tx.(*sql.Tx)

	/*
	err00 := torder.Get(1239, o)

	if err00 != nil {
		panic(err00)
	}
	*/
	findOpts1 := &TOrderFindOptions{}
	findOpts1.SetTx(tx)
	torders := make([]TOrder, 0)
	fmt.Println("findOpts: ", tx, err)
	err00 := torder.Find(&torders, findOpts1)

	fmt.Println("torderstorderstorderstorderstorders len: ", len(torders), fmt.Sprintf("%#v", torders))

	if err00 != nil {
		panic(err00)
	}
	fmt.Println("oooo: ", o.Id, o.OID, o.TradeNo)

	opts := options.Client().ApplyURI("mongodb://127.0.0.1:27017")

	c, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		fmt.Println("connect to mongodb failed")
		t.Error(err)
	}
	db := c.Database("guz-lib-test")

	es, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL("http://qa-new:9200/"),
	)
	if err != nil {
		fmt.Println("connect to elastic failed")
		t.Error("eeeeeeeeeeeeee", err)
	}

	dao2 := NewMTemplateDAO(db, es)

	tt := &ESTemplate{}
	err = dao2.e.Get("mcYdInkBIAk9r5vbgNCa", tt, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("tt.Name: ", tt.Name, "ppppp", tt.Id)
	os.Exit(0)
	tt.Name = "eee"
	tt.Tags = []string{fmt.Sprintf("%v", time.Now().Unix())}
	fmt.Println(*tt)

	id, err := dao2.e.Insert(tt, nil)
	fmt.Println("idididid", id, err)

	//mTemplate := &MTemplate{Template: tt.Template}
	//id, err3 := dao2.Insert(mTemplate, nil)
	//fmt.Println(id, err3)

	findOpts := NewMongodbFindOptions()
	data2 := make([]*MTemplate, 0)

	err = dao2.Find(&data2, findOpts)
	if err != nil {
		panic(err)
	}
	//for _, i := range data2 {
	//	fmt.Println("eeee", i.Bases)
	//}
	//fmt.Println(data2, findOpts.pagination.Total())

	fmt.Println("======================================================================")

}
