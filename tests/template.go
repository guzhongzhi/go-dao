package data

import (
	"github.com/guzhongzhi/gmicro/dao"
	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TemplateBase struct {
	Language      string `bson:"language" json:"language"`
	HomepageCover string `bson:"homepageCover" json:"homepageCover"`
	EditCover     string `bson:"editCover" json:"editCover"`
}

type Template struct {
	Name   string         `bson:"name" json:"name"`
	Tags   []string       `bson:"tags" json:"tags"`
	Status int            `bson:"status" json:"status"`
	Bases  []TemplateBase `bson:"bases" json:"bases"`
}

func NewMTemplateDAO(db *mongo.Database, es *elastic.Client) *templateDAO {
	return &templateDAO{
		dao.NewMongodbDAO(db, "e_template", nil),
		dao.NewElasticSearchDAO(es, "e_template"),
	}
}

type templateDAO struct {
	dao.MongodbDAO
	e dao.ElasticSearchDAO
}

func (s *templateDAO) DeleteById(id string) error {
	return s.e.Delete(id)
}

type FindOptions struct {
	id         *string
	tags       []string
	name       *string
	pagination *dao.Pagination
}

func (s *FindOptions) ID(v string) *FindOptions {
	s.id = &v
	return s
}

func (s *FindOptions) Tags(v []string) *FindOptions {
	s.tags = v
	return s
}

func (s *FindOptions) Name(v string) *FindOptions {
	s.name = &v
	return s
}

func (s *FindOptions) Pagination() *dao.Pagination {
	return s.pagination
}

func (s *FindOptions) Filter() (interface{}, error) {
	where := primitive.M{}
	if s.id != nil {
		objId, err := primitive.ObjectIDFromHex(*s.id)
		if err != nil {
			return nil, err
		}
		where["_id"] = objId
	}
	return where, nil
}

func (s *FindOptions) Options() (interface{}, error) {
	return nil, nil
}

type MongodbFindOptions FindOptions
type ElasticSearchFindOptions FindOptions

func NewMongodbFindOptions() *FindOptions {
	return &FindOptions{
		pagination: dao.NewPagination(),
	}
}
