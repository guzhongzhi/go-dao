package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongodbConfig struct {
	DSN string
}

func NewMongodbClient(cfg *MongodbConfig) *mongo.Client {
	opts := options.Client().ApplyURI(cfg.DSN)
	c, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		panic(err)
	}
	return c
}

type MongodbDAO interface {
	DAO
	Collection() *mongo.Collection
}

func NewMongodbDAO(db *mongo.Database, tableName string, opts options.CollectionOptions) MongodbDAO {
	return &mongodb{
		coll: db.Collection(tableName),
	}
}

func (s *mongodb) BeginTransaction(ctx context.Context, tx TxOptions) (interface{}, error) {
	return nil, fmt.Errorf("mongodb do not support transaction")
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

func (s *mongodb) FindOne(id interface{}, data Data, opts FindOptions) error {
	rs := s.coll.FindOne(context.Background(), primitive.M{"_id": id})
	return rs.Decode(data)
}

func (s *mongodb) Insert(entity Data, opts InsertOptions) (id interface{}, err error) {
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

func (s *mongodb) Update(id interface{}, data Data, updateOptions UpdateOptions) error {
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
