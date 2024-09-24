package mongo

import (
	"context"
	"errors"
	"time"

	libraryErrors "github.com/cristianat98/dbclientgo/errors"
	"github.com/cristianat98/dbclientgo/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoDb = "MongoDB"

// MongoManager is the structure to manage the connections and operations to the MongoDB
// client: It is directly the client to the MongoDB
// collection: It is the collection inside the MongoDB to make the queries. It comes from the client
// timeout: It is the time to wait before returns a TimeoutError
type MongoManager struct {
	client     *mongo.Client
	collection *mongo.Collection
	timeout    int64
}

// ConnectDB is the function inside the MongoManager to connect to the MongoDB
// dbUri: It is the URI to connect to the MongoDB
// dbName: It is the name of the DB inside the MongoDB
// collection: It is the name of the collection inside the DB inside the MongoDB
// timeout: It is the time to define the timeout inside the MongoManager
// It returns an error in case there was some error
func (manager *MongoManager) ConnectDB(dbUri, dbName, collection string, timeout int64) error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dbUri).SetServerAPIOptions(serverAPI)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	var result bson.M
	if err := client.Database(dbName).RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		logger.Error(err.Error())
		return &libraryErrors.ConnectionError{Db: mongoDb}
	}
	manager.client = client
	manager.collection = client.Database(dbName).Collection(collection)
	manager.timeout = timeout
	logger.Info("Connected to MongoDB")
	return nil
}

// DisconnectDB is the function inside the MongoManager to disconnect from the MongoDB
func (manager *MongoManager) DisconnectDB() {
	err := manager.client.Disconnect(context.TODO())
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Info("Disconnected to MongoDB")
}

// InsertOne is the function inside the MongoManager to Insert a document in the collection
// document: It is the document to add in the collection
// It returns the new document inserted in the collection and an error
func (manager *MongoManager) InsertOne(document map[string]interface{}) (map[string]interface{}, error) {
	ctx1, cancel1 := context.WithTimeout(context.Background(), time.Duration(manager.timeout)*time.Second)
	defer cancel1()

	resultInsert, err := manager.collection.InsertOne(ctx1, document)
	if err != nil {
		logger.Error(err.Error())
		switch err.(type) {
		case *mongo.CommandError:
			return nil, &libraryErrors.ConnectionError{Db: mongoDb}
		case *mongo.WriteErrors:
			return nil, &libraryErrors.AlreadyExistError{Message: "Document with a key already exists"}
		default:
			return nil, err
		}
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Duration(manager.timeout)*time.Second)
	defer cancel2()
	resultFind := manager.collection.FindOne(ctx2, map[string]interface{}{"_id": resultInsert.InsertedID})
	if err := resultFind.Err(); err != nil {
		logger.Error(err.Error())
		switch err.(type) {
		case *mongo.CommandError:
			return nil, &libraryErrors.ConnectionError{Db: mongoDb}
		default:
			return nil, err
		}
	}
	var documentReturned bson.M
	err = resultFind.Decode(&documentReturned)
	if err != nil {
		logger.Error(err.Error())
	}
	return documentReturned, err
}

// Find is the function inside the MongoManager to return a list of documents that match the filter defined
// filter: It is the filter to find documents inside the MongoDB
// It returns a list of documents (it may be empty) and an error
func (manager *MongoManager) Find(filter map[string]interface{}) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(manager.timeout)*time.Second)
	defer cancel()

	var results []map[string]interface{}
	cursor, err := manager.collection.Find(ctx, filter)
	if err != nil {
		logger.Error(err.Error())
		switch err.(type) {
		case *mongo.CommandError:
			return nil, &libraryErrors.ConnectionError{Db: mongoDb}
		default:
			return nil, err
		}
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var result map[string]interface{}

		err = cursor.Decode(&result)
		results = append(results, result)
	}

	return results, err
}

// UpdateOne is the function inside the MongoManager to update the first document that matches with the filter defined
// filter: It is the filter to find documents inside the MongoDB to update
// document: It contains the new changes to apply in the document found
// It returns the document updated and an error
func (manager *MongoManager) UpdateOne(filter map[string]interface{}, document interface{}) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(manager.timeout)*time.Second)
	defer cancel()

	resultFind, err := manager.Find(filter)
	if err != nil {
		return nil, err
	}
	if len(resultFind) != 1 {
		return nil, errors.New("document not found")
	}
	id := resultFind[0]["_id"]
	filterInternal := map[string]interface{}{
		"_id": id,
	}
	m, ok := document.(map[string]interface{})
	if !ok {
		logger.Error("i is not a map[string]interface{}")
		return nil, nil
	}

	_, err = manager.collection.UpdateOne(ctx, filterInternal, bson.M{"$set": m})
	if err != nil {
		return nil, err
	}
	resultFind, _ = manager.Find(filterInternal)
	return resultFind[0], nil
}

// DeleteMany is the function inside the MongoManager to delete all the documents that match with the filter
// It returns the number of documents deleted and an error
func (manager *MongoManager) DeleteMany(filter map[string]interface{}) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(manager.timeout)*time.Second)
	defer cancel()

	result, err := manager.collection.DeleteMany(ctx, filter)
	return int(result.DeletedCount), err
}
