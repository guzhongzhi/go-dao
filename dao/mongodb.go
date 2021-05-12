package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongodbDAO interface {
	DAO
	Collection() *mongo.Collection
}

func NewMongodbDAO(db *mongo.Database, tableName string, opts options.CollectionOptions) MongodbDAO {
	return &mongodb{
		coll: db.Collection(tableName),
	}
}

type mongodb struct {
	coll *mongo.Collection
}

func (s *mongodb) Collection() *mongo.Collection {
	return s.coll
}

func (s *mongodb) Delete(id interface{}, opts DeleteOptions) error {
	_, err := s.coll.DeleteOne(context.Background(), primitive.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (s *mongodb) Get(id interface{}, data Entity, opts GetOptions) error {
	rs := s.coll.FindOne(context.Background(), primitive.M{"_id": id})
	return rs.Decode(data)
}

func (s *mongodb) Insert(entity Entity, opts InsertOptions) (id interface{}, err error) {
	if !entity.IsNew() {
		return nil, fmt.Errorf("collection '%s' insert error, the data is not a new record", s.coll.Name())
	}
	entity.SetID(primitive.NewObjectID())
	rs, err := s.coll.InsertOne(context.Background(), entity)
	if err != nil {
		return nil, fmt.Errorf("collection '%s' insert error: '%s'", s.coll.Name(), err.Error())
	}
	entity.SetID(rs.InsertedID)
	return rs.InsertedID, err
}

func (s *mongodb) Update(id interface{}, data Entity, updateOptions UpdateOptions) error {
	_, err := s.coll.UpdateOne(context.Background(), primitive.M{
		"_id": id,
	}, primitive.M{"$set": data})
	if err != nil {
		return fmt.Errorf("collection '%s' update error: '%s'", s.coll.Name(), err.Error())
	}
	return nil
}

func (s *mongodb) Find(data interface{}, opts FindOptions) error {

	findOpts := &options.FindOptions{}
	if opts.Pagination() != nil {
		findOpts.SetLimit(opts.Pagination().PageSize()).SetSkip(opts.Pagination().Offset())
	}

	filter, err := opts.Filter()
	if err != nil {
		return err
	}

	total, err := s.coll.CountDocuments(context.Background(), filter)
	if err != nil {
		return err
	}
	opts.Pagination().SetTotal(total)

	cursor, err := s.coll.Find(context.Background(), filter, findOpts)
	if err != nil {
		return fmt.Errorf("collection '%s' find operation error: '%s'", s.coll.Name(), err.Error())
	}
	err = cursor.All(context.Background(), data)
	fmt.Println("data: ", data)
	if err != nil {
		return fmt.Errorf("collection '%s' cursor all operation error: '%s'", s.coll.Name(), err.Error())
	}
	return nil
}
