package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/pinguo-icc/salad-effect/internal/infrastructure/repository"
	"reflect"
)

type ElasticSearchDAO interface {
	DAO
	Client() *elastic.Client
	IndexName() string
}

func NewElasticSearchDAO(client *elastic.Client, indexName string) ElasticSearchDAO {
	e := &elasticsearch{
		client: client,
		index:  indexName,
	}
	e.init()
	return e
}

type elasticsearch struct {
	client *elastic.Client
	index  string
}

func (s *elasticsearch) IndexName() string {
	return s.index
}

func (s *elasticsearch) Client() *elastic.Client {
	return s.client
}

func (s *elasticsearch) init() {
	b, _ := s.client.IndexExists(s.index).Do(context.Background())
	if !b {
		s.client.CreateIndex(s.index).Do(context.Background())
	}
}

func (s *elasticsearch) Insert(entity Entity, opts InsertOptions) (id interface{}, err error) {
	js, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}
	rsp, err := s.client.Index().Index(s.index).BodyJson(string(js)).Do(context.Background())
	if err != nil {
		return nil, err
	}
	entity.SetID(rsp.Id)
	return rsp.Id, nil
}

func (s *elasticsearch) Update(id interface{}, data Entity, opts UpdateOptions) error {
	_, err := s.client.Update().
		Index(s.index).
		Id(fmt.Sprintf("%s", id)).
		Doc(data).
		Do(context.Background())

	return err
}

func (s *elasticsearch) Find(opts interface{}, data interface{}) error {
	rv := reflect.ValueOf(data)
	if rv.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("elasticsearch find data should be pointer of struct or map for the index of '%s'", s.index)
	}
	rv = rv.Elem()
	if !rv.CanSet() {
		return fmt.Errorf("elasticsearch find data can not be set for index of '%s'", s.index)
	}

	o := opts.(FindOptions)
	filter, err := o.Filter()
	if err != nil {
		return err
	}
	query := s.client.Search().
		Index(s.index).
		Query(filter.(elastic.Query))

	if o.Pagination() != nil {
		query.Size(int(o.Pagination().PageSize())).
			From(int(o.Pagination().Offset()))
	}

	rs, err := query.Do(context.Background())

	if err != nil {
		return err
	}

	isSlice := rv.Type().Kind() == reflect.Slice
	item := rv.Addr().Interface()
	if isSlice {
		item = reflect.New(rv.Type().Elem()).Interface()
	}

	if !isSlice && rs.Hits.TotalHits.Value == 0 {
		return repository.ErrNotFound
	}

	for _, it := range rs.Hits.Hits {
		if err := json.Unmarshal(it.Source, item); err != nil {
			return err
		}
		if !isSlice {
			rv.Set(reflect.ValueOf(item).Elem())
			return nil
		}
		rv.Set(reflect.Append(rv, reflect.ValueOf(item).Elem()))
	}
	return nil
}

func (s *elasticsearch) Delete(id interface{}) error {
	_, err := s.client.Delete().
		Index(s.index).
		Id(fmt.Sprintf("%s", id)).
		Do(context.Background())

	if err != nil {
		return err
	}
	return nil
}

func (s *elasticsearch) Get(id interface{}, data Entity) error {
	rs, err := s.client.Get().
		Index(s.index).
		Id(fmt.Sprintf("%s", id)).
		Do(context.Background())

	if err != nil {
		return err
	}
	body, err := rs.Source.MarshalJSON()
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}
	data.SetID(rs.Id)

	return nil
}
