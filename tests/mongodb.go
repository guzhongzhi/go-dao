package data

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MTemplate struct {
	Id       primitive.ObjectID `bson:"_id"`
	Template `bson:",inline"`
}

func (s *MTemplate) SetID(v interface{}) {
	s.Id = v.(primitive.ObjectID)
}

func (s *MTemplate) ID() interface{} {
	return s.Id
}

func (s *MTemplate) IsNew() bool {
	return s.Id.IsZero()
}

func (s *MTemplate) String() string {
	return s.Id.Hex()
}
