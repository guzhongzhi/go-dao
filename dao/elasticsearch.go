package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/guzhongzhi/gmicro/logger"
	"github.com/olivere/elastic/v7"
	"reflect"
)

type ElasticSearchOptions struct {
	Mapping   string
	IndexName string
}

type ElasticSearchDAO interface {
	DAO
	Client() *elastic.Client
	IndexName() string
}

type ElasticSearchConfig struct {
	Addresses   []string
	Sniff       bool
	HealthCheck bool
	Username    string
	Password    string
}

func (s *ElasticSearchConfig) toOptions() []elastic.ClientOptionFunc {
	opts := make([]elastic.ClientOptionFunc, 0)

	opts = append(opts, elastic.SetSniff(s.Sniff))
	opts = append(opts, elastic.SetURL(s.Addresses...))
	opts = append(opts, elastic.SetHealthcheck(s.HealthCheck))
	opts = append(opts, elastic.SetSniff(s.Sniff))
	if s.Username != "" && s.Password != "" {
		opts = append(opts, elastic.SetBasicAuth(s.Username, s.Password))
	}
	return opts
}

func NewElasticClient(cfg *ElasticSearchConfig) *elastic.Client {
	es, err := elastic.NewClient(cfg.toOptions()...)
	if err != nil {
		panic(err)
	}
	return es
}

func NewElasticSearchDAO(client *elastic.Client, options ElasticSearchOptions, log logger.SuperLogger) ElasticSearchDAO {
	if log == nil {
		log = logger.Default()
	}
	e := &elasticsearch{
		client:  client,
		index:   options.IndexName,
		options: options,
		logger:  log,
	}
	e.init()
	return e
}

type elasticsearch struct {
	client  *elastic.Client
	index   string
	options ElasticSearchOptions
	logger  logger.SuperLogger
}

func (s *elasticsearch) BeginTransaction(ctx context.Context, tx TxOptions) (interface{}, error) {
	return nil, fmt.Errorf("elasticsearch do not support transaction")
}

func (s *elasticsearch) IndexName() string {
	return s.index
}

func (s *elasticsearch) Client() *elastic.Client {
	return s.client
}

func (s *elasticsearch) init() {
	b, err := s.client.IndexExists(s.index).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if !b {
		icr, err := s.client.CreateIndex(s.index).BodyString(s.options.Mapping).Do(context.Background())
		if err != nil {
			panic(err)
		}
		s.logger.Debugf("create new index '%s' succeed,%v", s.index, icr.Acknowledged)
	}
}

func (s *elasticsearch) Insert(entity Data, opts InsertOptions) (id interface{}, err error) {
	js, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}
	q := s.client.Index().Index(s.index).BodyJson(string(js))
	if !entity.IsNew() {
		q.Id(fmt.Sprintf("%v", entity.GetID()))
	}
	rsp, err := q.Do(context.Background())
	if err != nil {
		return nil, err
	}
	entity.SetID(rsp.Id)
	return rsp.Id, nil
}

func (s *elasticsearch) Update(id interface{}, data Data, opts UpdateOptions) error {
	s.logger.Debugf("update doc '%s'.'%s'", s.index, id)
	updateService := s.client.Update().
		Index(s.index).
		Id(fmt.Sprintf("%s", id)).
		Doc(data)

	_, err := updateService.Do(context.Background())

	return err
}

func (s *elasticsearch) Find(data interface{}, opts FindOptions) error {
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
		return nil
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

func (s *elasticsearch) Delete(id interface{}, opts DeleteOptions) error {
	_, err := s.client.Delete().
		Index(s.index).
		Id(fmt.Sprintf("%s", id)).
		Do(context.Background())

	if err != nil {
		return err
	}
	return nil
}

func (s *elasticsearch) FindOne(id interface{}, data Data, opts FindOptions) error {
	rs, err := s.client.Get().
		Index(s.index).
		Id(fmt.Sprintf("%s", id)).
		Do(context.Background())

	if err != nil {
		if !elastic.IsNotFound(err) {
			return err
		} else {
			return nil
		}
	}

	err = json.Unmarshal(rs.Source, data)
	data.SetID(rs.Id)

	return nil
}
