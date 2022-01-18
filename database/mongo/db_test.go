package mongo

import (
	"context"
	"flag"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/power28-china/auth/config"
	"github.com/power28-china/auth/utils/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type document struct {
	Name   string
	Gender string
	Role   string
}

var (
	documents  []document
	collection *mongo.Collection
	faker      = gofakeit.New(0)
	cursor     *mongo.Cursor
)

const collectionName = "test"

func TestMain(m *testing.M) {
	flag.Parse()

	for i := 1; i <= 100; i++ {
		d := document{faker.Name(), faker.Gender(), faker.JobTitle()}
		documents = append(documents, d)
	}

	exitCode := m.Run()

	documents = nil

	os.Exit(exitCode)
}

func TestSave(t *testing.T) {
	result := Save(collectionName, document{faker.Name(), faker.Gender(), faker.JobTitle()})
	if result.InsertedID == nil {
		t.Errorf("Save test failed.")
	}
}

func TestSaveAll(t *testing.T) {
	var interfaceSlice []interface{} = make([]interface{}, len(documents))
	for i, d := range documents {
		interfaceSlice[i] = d
	}
	result := SaveAll(collectionName, interfaceSlice)
	if len(result.InsertedIDs) != len(documents) {
		t.Errorf("SaveAll test failed.")
	}
	for _, id := range result.InsertedIDs {
		logger.Sugar.Debug(id)
	}
}

func TestDelete(t *testing.T) {
	d := document(documents[0])
	result := Delete(collectionName, "name", d.Name)
	if result != 1 {
		t.Errorf("Delete test failed.")
	}
}

func TestDeleteAll(t *testing.T) {
	d := document(documents[2])
	filter := bson.D{primitive.E{Key: "name", Value: d.Name}}
	result := DeleteAll(collectionName, filter)
	if result < 1 {
		t.Errorf("DeleteAll test failed.")
	}
}

func TestUpdate(t *testing.T) {
	d := document(documents[5])
	filter := bson.D{primitive.E{Key: "name", Value: d.Name}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "role", Value: "employee"}}}}
	result := Update(collectionName, filter, update)
	if result != 1 {
		t.Errorf("Update test failed.")
	}
}

func TestUpdateAll(t *testing.T) {
	d := document(documents[3])
	filter := bson.D{primitive.E{Key: "role", Value: d.Role}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "role", Value: "employee"}}}}
	result := UpdateAll(collectionName, filter, update)
	if result < 1 {
		t.Errorf("DeleteAll test failed.")
	}
}

func TestFind(t *testing.T) {
	d := document(documents[3])
	res := &document{}
	if err := Find(collectionName, "name", d.Name).Decode(res); err != nil {
		t.Errorf("Find test failed.")
	}

	logger.Sugar.Debug("Find one document: ", res)
}

func TestFindAll(t *testing.T) {
	res := []document{}
	filter := bson.D{primitive.E{Key: "gender", Value: "male"}}
	if cursor, err = FindAll(collectionName, filter); err != nil {
		t.Errorf("FindAll test failed.")
	}

	//延迟关闭游标
	defer func() {
		if err = cursor.Close(context.TODO()); err != nil {
			logger.Sugar.Fatal(err)
		}
	}()

	if err = cursor.All(context.TODO(), &res); err != nil {
		logger.Sugar.Fatal(err)
	}
	// for _, result := range res {
	// 	logger.Sugar.Debug(result)
	// }

	logger.Sugar.Debugf("Find %d documents.\n", len(res))
}

func TestFindAllByFilters(t *testing.T) {
	res := []document{}
	filter := bson.D{primitive.E{Key: "gender", Value: "female"}}
	if cursor, err = FindAll(collectionName, filter); err != nil {
		t.Errorf("FindAll test failed.")
	}

	//延迟关闭游标
	defer func() {
		if err = cursor.Close(context.TODO()); err != nil {
			logger.Sugar.Fatal(err)
		}
	}()

	if err = cursor.All(context.TODO(), &res); err != nil {
		logger.Sugar.Fatal(err)
	}
	// for _, result := range res {
	// 	logger.Sugar.Debug(result)
	// }

	logger.Sugar.Debugf("Find %d documents.\n", len(res))
}

func TestCollectionCount(t *testing.T) {
	name, count := CollectionCount(collectionName, bson.D{})
	if (name == "") || (count == 0) {
		t.Errorf("CollectionCount test failed.\n")
	}
	logger.Sugar.Debugf("Collection name:%s, %d documents in it.\n", name, count)
}

func TestFindWithOptions(t *testing.T) {
	var results []bson.M

	projection := bson.D{
		{"_id", 0},
		{"provinces", bson.D{
			{"$elemMatch", bson.D{
				{"label", "湖北省"},
			}},
		}},
	}

	options := options.Find().SetProjection(projection)
	cursor, err := FindWithOptions(config.Config("COUNTRY_COLLECTION"), bson.D{}, options)
	if err != nil {
		t.Errorf("FindWithOptions failed: %#v\n", err)
	}
	if err := cursor.All(context.TODO(), &results); err != nil {
		t.Errorf("FindWithOptions failed: %v\n", err)
	}
	logger.Sugar.Debugf("Results: %v\n", results)
}
