package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"reflect"
	"regexp"
	"strings"
)

//type basic[T any] struct{}
//

var matchFirstCapSnakeCase = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCapSnakeCase = regexp.MustCompile("([a-z\\d])([A-Z])")

// SnakeCase is
func snakeCase(str string) string {
	snake := matchFirstCapSnakeCase.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCapSnakeCase.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func toSliceAny[T any](objs []T) []any {
	var results []any
	for _, obj := range objs {
		results = append(results, obj)
	}
	return results
}

//type AdapterGateway[T any] struct{}
//
//func (g *AdapterGateway[T]) InsertOrUpdate(obj *T) error {
//	return nil
//}
//
//func (g *AdapterGateway[T]) InsertMany(objs ...*T) error {
//	return nil
//}
//
//func (g *AdapterGateway[T]) GetOne(filter map[string]any, result *T) error {
//	return nil
//}
//
//func (g *AdapterGateway[T]) GetAll(param GetAllParam, results *[]*T) (int64, error) {
//	return 0, nil
//}
//
//func (g *AdapterGateway[T]) GetAllEachItem(param GetAllParam, resultEachItem func(result T)) (int64, error) {
//	return 0, nil
//}

// =======================================

type MongoGateway[T any] struct {
	Database *mongo.Database
}

func NewMongoGateway[T any](db *mongo.Database) *MongoGateway[T] {
	return &MongoGateway[T]{
		Database: db,
	}
}

func NewDatabase(databaseName string) *mongo.Database {

	uri := "mongodb://localhost:27017/?readPreference=primary&ssl=false"

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))

	err = client.Connect(context.Background())
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		panic(err)
	}

	return client.Database(databaseName)

}

func (g *MongoGateway[T]) GetTypeName() string {
	var x T
	return snakeCase(reflect.TypeOf(x).Name())
}

//func (g *MongoGateway[T]) GetCollection() *mongo.Collection {
//	var x T
//	name := snakeCase(reflect.TypeOf(x).Name())
//	return g.Database.Collection(name)
//}

func (g *MongoGateway[T]) InsertOrUpdate(obj *T) error {

	sf, exist := reflect.TypeOf(obj).Elem().FieldByName("ID")
	if !exist {
		return fmt.Errorf("field ID as primary key is not found in %s", reflect.TypeOf(obj).Name())
	}

	tagValue, exist := sf.Tag.Lookup("bson")
	if !exist || tagValue != "_id" {
		return fmt.Errorf("field ID must have tag `bson:\"_id\"`")
	}

	filter := bson.D{{"_id", reflect.ValueOf(obj).Elem().FieldByName("ID").Interface()}}
	update := bson.D{{"$set", obj}}
	opts := options.Update().SetUpsert(true)

	coll := g.Database.Collection(g.GetTypeName())
	_, err := coll.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func (g *MongoGateway[T]) InsertMany(objs ...*T) error {

	if len(objs) == 0 {
		return fmt.Errorf("objs must > 0")
	}

	opts := options.InsertMany().SetOrdered(false)

	coll := g.Database.Collection(g.GetTypeName())
	_, err := coll.InsertMany(context.TODO(), toSliceAny(objs), opts)
	if err != nil {
		return err
	}

	return nil
}

func (g *MongoGateway[T]) GetOne(filter map[string]any, result *T) error {

	coll := g.Database.Collection(g.GetTypeName())

	singleResult := coll.FindOne(context.TODO(), filter)

	err := singleResult.Decode(result)
	if err != nil {
		return err
	}

	return nil
}

func (g *MongoGateway[T]) GetAll(param GetAllParam, results *[]*T) (int64, error) {

	coll := g.Database.Collection(g.GetTypeName())

	skip := param.Size * (param.Page - 1)
	limit := param.Size

	findOpts := options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
		Sort:  param.Sort,
	}

	ctx := context.TODO()

	count, err := coll.CountDocuments(ctx, param.Filter)
	if err != nil {
		return 0, err
	}

	cursor, err := coll.Find(ctx, param.Filter, &findOpts)
	if err != nil {
		return 0, err
	}

	err = cursor.All(ctx, results)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (g *MongoGateway[T]) GetAllEachItem(param GetAllParam, resultEachItem func(result T)) (int64, error) {

	coll := g.Database.Collection(g.GetTypeName())

	skip := param.Size * (param.Page - 1)
	limit := param.Size

	findOpts := options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
		Sort:  param.Sort,
	}

	ctx := context.TODO()

	count, err := coll.CountDocuments(ctx, param.Filter)
	if err != nil {
		return 0, err
	}

	cursor, err := coll.Find(ctx, param.Filter, &findOpts)
	if err != nil {
		return 0, err
	}

	for cursor.Next(ctx) {

		var result T
		err := cursor.Decode(&result)
		if err != nil {
			return 0, err
		}

		resultEachItem(result)

	}

	err = cursor.Err()
	if err != nil {
		return 0, err
	}

	return count, nil

}

func (g *MongoGateway[T]) Delete(filter map[string]any) error {

	coll := g.Database.Collection(g.GetTypeName())

	_, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}
