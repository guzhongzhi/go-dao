package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongodbConfig struct {
	DSN      string
	Database string
}

func NewMongodbDatabase(cfg *MongodbConfig) *mongo.Database {
	client := NewMongodbClient(cfg)
	fmt.Println(cfg.DSN)
	return client.Database(cfg.Database)
}

func NewMongodbClient(cfg *MongodbConfig) *mongo.Client {
	opts := options.Client().ApplyURI(cfg.DSN)
	c, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		panic(err)
	}
	return c
}

type MongodbIndex struct {
	Name   string
	Unique bool
}

type MongodbDAO interface {
	DAO
	Collection() *mongo.Collection
	CreateIndex(name string, keys interface{}, indexOptions *options.IndexOptions) error
	NewIndexOptions() *options.IndexOptions
	Indexes() []*MongodbIndex
	DropIndex(name string) error
	InsertMany(entities []interface{}, opts InsertOptions) (ids []interface{}, err error)
}

func NewMongodbDAO(db *mongo.Database, tableName string, opts ...*options.CollectionOptions) MongodbDAO {
	v := &mongodb{
		coll: db.Collection(tableName, opts...),
	}
	v.init()
	return v
}

func (s *mongodb) BeginTransaction(ctx context.Context, tx TxOptions) (interface{}, error) {
	return nil, fmt.Errorf("mongodb do not support transaction")
}

type mongodb struct {
	coll    *mongo.Collection
	indexes []*MongodbIndex
}

func (s *mongodb) init() {
	ctx := context.Background()
	cursor, err := s.Collection().Indexes().List(ctx)
	if err != nil {
		panic(err)
	}
	indexes := make([]*MongodbIndex, 0)
	err = cursor.All(context.Background(), &indexes)
	if err != nil {
		panic(err)
	}
	s.indexes = indexes
}

func (s *mongodb) Indexes() []*MongodbIndex {
	return s.indexes
}

func (s *mongodb) newIndexOptions() *options.IndexOptions {
	opts := &options.IndexOptions{}
	return opts
}

func (s *mongodb) NewIndexOptions() *options.IndexOptions {
	opts := &options.IndexOptions{}
	return opts
}

func (s *mongodb) DropIndex(name string) error {
	_, err := s.coll.Indexes().DropOne(context.Background(), name)
	if err != nil {
		return err
	}
	s.init()
	return nil
}

func (s *mongodb) CreateIndex(name string, keys interface{}, indexOptions *options.IndexOptions) error {
	for _, idx := range s.indexes {
		if idx.Name == name {
			return nil
		}
	}
	indexOptions.SetName(name)
	ctx := context.Background()
	indexModel := mongo.IndexModel{
		Keys:    keys,
		Options: indexOptions,
	}
	_, err := s.Collection().Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		panic(err)
	}
	s.init()
	return nil
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

func (s *mongodb) Count(opts FindOptions) (int64, error) {
	filter, err := opts.Filter()
	if err != nil {
		return 0, err
	}
	return s.coll.CountDocuments(context.Background(), filter)
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

	if opts.Sorts() != nil {
		findOpts.SetSort(opts.Sorts())
	}

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

func (s *mongodb) InsertMany(entities []interface{}, opts InsertOptions) (ids []interface{}, err error) {

	rs, err := s.coll.InsertMany(context.Background(), entities)
	if err != nil {
		return nil, fmt.Errorf("collection '%s' insert error: '%s'", s.coll.Name(), err.Error())
	}
	return rs.InsertedIDs, err
}
