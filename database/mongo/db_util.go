package mongo

import (
	"context"
	"time"

	"github.com/power28-china/auth/utils"
	"github.com/power28-china/auth/utils/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type index struct {
	Key  bson.D
	Name string
}

// Matchby represents the options for a matching.
type Matchby struct {
	Key   string
	Value string
}

// Save saves one object to the database.
func Save(c string, document interface{}) *mongo.InsertOneResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := db.Collection(c).InsertOne(ctx, document)
	if err != nil {
		logger.Sugar.Fatal(err)
		return nil
	}
	// logger.Sugar.Debugf("Inserted a single document %s successfully\n", result.InsertedID)

	return result
}

// SaveAll save many objects to the database.
func SaveAll(c string, documents []interface{}) *mongo.InsertManyResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := db.Collection(c).InsertMany(ctx, documents)
	if err != nil {
		logger.Sugar.Fatal(err)
	}
	logger.Sugar.Debugf("Inserted %d documents successfully\n", len(result.InsertedIDs))

	return result
}

// Delete deletes a single document.
func Delete(c string, key string, document interface{}) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{key: document}
	count, err := db.Collection(c).DeleteOne(ctx, filter, nil)
	if err != nil {
		logger.Sugar.Fatal(err)
	}
	logger.Sugar.Debugf("Delete %d document successfully\n", count.DeletedCount)

	return count.DeletedCount
}

// DeleteAll deletes all documents
func DeleteAll(c string, filter interface{}) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := db.Collection(c).DeleteMany(ctx, filter)
	if err != nil {
		logger.Sugar.Fatal(err)
	}
	logger.Sugar.Debugf("Delete %d documents successfully\n", count.DeletedCount)

	return count.DeletedCount
}

// Replace replace a document with a new one
func Replace(c string, filter, documents interface{}) int64 {
	// ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	// defer cancel()

	// specify the Upsert option to insert a new document if a document matching the filter isn't found
	opts := options.Replace().SetUpsert(true)

	result, err := db.Collection(c).ReplaceOne(context.Background(), filter, documents, opts)
	if err != nil {
		logger.Sugar.Fatal(err)
	}

	if result.ModifiedCount != 0 {
		// logger.Sugar.Debug("matched and replaced an existing document")
		return result.ModifiedCount
	}

	// logger.Sugar.Debugf("inserted a new document with ID %v\n", result.UpsertedID)
	return result.UpsertedCount
}

// Update update or insert document.
func Update(c string, filter, document interface{}) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// specify the Upsert option to insert a new document if a document matching the filter isn't found
	opts := options.Update().SetUpsert(true)
	result, err := db.Collection(c).UpdateOne(ctx, filter, document, opts)
	if err != nil {
		logger.Sugar.Fatal(err)
	}

	if result.ModifiedCount != 0 {
		logger.Sugar.Debugf("Update %d documents successfully\n", result.ModifiedCount)
		return result.ModifiedCount
	}
	logger.Sugar.Debugf("Update documents %s successfully\n", result.UpsertedID)
	return result.UpsertedCount
}

// UpdateAll update all documents in the collection
func UpdateAll(c string, filter, documents interface{}) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := db.Collection(c).UpdateMany(ctx, filter, documents)
	if err != nil {
		logger.Sugar.Fatal(err)
	}

	if result.ModifiedCount != 0 {
		logger.Sugar.Debugf("Update %d documents successfully\n", result.ModifiedCount)
		return result.ModifiedCount
	}
	logger.Sugar.Debugf("Update documents %s successfully\n", result.UpsertedID)
	return result.UpsertedCount
}

// Find returns a document matching the given filter in the collection.
func Find(c string, key string, document interface{}) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection, err := db.Collection(c).Clone()
	if err != nil {
		logger.Sugar.Fatal(err)
	}

	filter := bson.D{primitive.E{Key: key, Value: document}}
	singleResult := collection.FindOne(ctx, filter)
	return singleResult
}

// FindOneAndUpdate updates a single document based on the filter with the document.
func FindOneAndUpdate(c string, filter interface{}, document interface{}) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection, err := db.Collection(c).Clone()
	if err != nil {
		logger.Sugar.Fatal(err)
	}

	opts := options.FindOneAndUpdate().SetUpsert(true)

	singleResult := collection.FindOneAndUpdate(ctx, filter, document, opts)
	return singleResult
}

// FindAll returns many documents matching the filter in the collection.
func FindAll(c string, filter interface{}) (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection, err := db.Collection(c).Clone()
	if err != nil {
		logger.Sugar.Fatal(err)
	}
	return collection.Find(ctx, filter)
}

// FindWithOptions returns matched documents with find options.
func FindWithOptions(c string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection, err := db.Collection(c).Clone()
	if err != nil {
		logger.Sugar.Fatal(err)
	}

	return collection.Find(ctx, filter, opts...)
}

// CollectionCount return how many documents in the collection.
func CollectionCount(c string, filters interface{}, opts ...*options.CountOptions) (string, int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(c)
	name := collection.Name()
	size, _ := collection.CountDocuments(ctx, filters, opts...)
	return name, size
}

// CollectionDocuments returns many documents matching the given filter and other criteria options in the collection. sort 1 ascending and -1 descending.
func CollectionDocuments(c string, Skip, Limit int64, sort, filter interface{}) *mongo.Cursor {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find().SetSort(sort).SetLimit(Limit).SetSkip(Skip)
	temp, err := db.Collection(c).Find(ctx, filter, findOptions)
	if err != nil {
		logger.Sugar.Fatal(err)
	}
	return temp
}

// CreateTTLIndex Create TTL Index for specified collection. returns the index name and error if it occurs.
func CreateTTLIndex(c string, expireAfterSeconds int32) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if !isIndexExist(db.Collection(c).Indexes(), index{Name: "expiresin_1"}) {
		indexModel := mongo.IndexModel{
			Keys:    bson.D{primitive.E{Key: "expiresin", Value: 1}}, // index in ascending order.
			Options: options.Index().SetBackground(true).SetExpireAfterSeconds(expireAfterSeconds),
		}

		name, err := db.Collection(c).Indexes().CreateOne(ctx, indexModel)
		if err != nil {
			return "", err
		}
		return name, nil
	}

	return "", nil
}

func isIndexExist(iv mongo.IndexView, expected index) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := iv.List(ctx)
	if err != nil {
		logger.Sugar.Fatal("List error: ", err)
	}

	found := false
	for cursor.Next(ctx) {
		var idx index
		err = cursor.Decode(&idx)
		if err != nil {
			logger.Sugar.Fatal("Decode error: ", err)
		}

		if idx.Name == expected.Name {
			found = true
		}
	}
	return found
}

// Sum returns the sum of all the accounts for given key.
func Sum(c string, pipeline string) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var cursor *mongo.Cursor

	opts := options.Aggregate().SetMaxTime(5 * time.Second)

	cursor, err = db.Collection(c).Aggregate(ctx, utils.MongoPipeline(pipeline), opts)

	logger.Sugar.Debugf("Cursor: %#v", cursor)

	if err != nil {
		return nil, err
	}

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}
