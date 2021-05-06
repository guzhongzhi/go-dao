package data

import (
	"github.com/olivere/elastic/v7"
)

type ESTemplate struct {
	Template `json:",inline"`
	Id       string `json:"id"`
}

func (s *ESTemplate) ID() interface{} {
	return s.Id
}

func (s *ESTemplate) SetID(v interface{}) {
	s.Id = v.(string)
}

func (s *ESTemplate) IsNew() bool {
	return s.Id != ""
}

func (s *ESTemplate) String() string {
	return s.Id
}

type ESFindOptions FindOptions

func (s *ESFindOptions) Filter() interface{} {
	return elastic.NewMatchAllQuery()
}
