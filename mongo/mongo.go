package mongo

import (
	"context"
	"errors"
	"reflect"
	"time"

	libraryErrors "github.com/cristianat98/dbclientgo/errors"
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

// CreateMongoManager is the constructor for the MongoManager. If it can not connect to the MongoDB, it will fail
// dbUri: It is the URI to connect to the MongoDB
// dbName: It is the name of the DB inside the MongoDB
// collection: It is the name of the collection inside the DB inside the MongoDB
// timeout: It is the time to define the timeout inside the MongoManager
// It returns the MongoManager instance and an error
func CreateMongoManager(dbUri, dbName, collection string, timeout int64) (*MongoManager, error) {
	var mongoManager = new(MongoManager)
	if err := mongoManager.ConnectDb(dbUri, dbName, collection, timeout); err != nil {
		return nil, err
	}
	return mongoManager, nil
}

// ConnectDb is the function inside the MongoManager to connect to the MongoDB
// dbUri: It is the URI to connect to the MongoDB
// dbName: It is the name of the DB inside the MongoDB
// collection: It is the name of the collection inside the DB inside the MongoDB
// timeout: It is the time to define the timeout inside the MongoManager
// It returns an error in case there was some error
func (manager *MongoManager) ConnectDb(dbUri, dbName, collection string, timeout int64) error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dbUri).SetServerAPIOptions(serverAPI)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}

	var result bson.M
	if err := client.Database(dbName).RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return &libraryErrors.ConnectionError{Db: mongoDb}
	}
	manager.client = client
	manager.collection = client.Database(dbName).Collection(collection)
	manager.timeout = timeout
	return nil
}

// DisconnectDb is the function inside the MongoManager to disconnect from the MongoDB
func (manager *MongoManager) DisconnectDb() {
	manager.client.Disconnect(context.TODO())
}

// InsertOne is the function inside the MongoManager to insert a document in the collection
// document: It is the document to add in the collection
// It returns the new document inserted in the collection and an error
func (manager *MongoManager) InsertOne(document map[string]interface{}) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(manager.timeout)*time.Second)
	defer cancel()

	resultInsert, err := manager.collection.InsertOne(ctx, document)
	if err != nil {
		if _, ok := err.(mongo.CommandError); ok {
			return nil, &libraryErrors.ConnectionError{Db: mongoDb}
		} else if _, ok := err.(mongo.WriteException); ok {
			return nil, &libraryErrors.AlreadyExistError{Message: "Document with a key already exists"}
		} else {
			return nil, err
		}
	}

	documentReturned, err := manager.FindOne(map[string]interface{}{"_id": resultInsert.InsertedID})
	return documentReturned, err
}

// InsertMany is the function inside the MongoManager to insert many documents in the collection
// documents: It is the list of documents to insert in the collection
// It returns the new documents inserted in the collection and an error
func (manager *MongoManager) InsertMany(documents []map[string]interface{}) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(manager.timeout)*time.Second)
	defer cancel()

	var documentsParsed []interface{}
	for _, item := range documents {
		documentsParsed = append(documentsParsed, item)
	}

	insertResult, err := manager.collection.InsertMany(ctx, documentsParsed)
	var documentsInserted []map[string]interface{}

	if err != nil {
		for _, id := range insertResult.InsertedIDs {
			documentReturned, err := manager.FindOne(map[string]interface{}{"_id": id})
			if err != nil {
				if _, ok := err.(mongo.CommandError); ok {
					return nil, &libraryErrors.ConnectionError{Db: mongoDb}
				} else {
					return nil, err
				}
			}
			documentsInserted = append(documentsInserted, documentReturned)
		}
		if _, ok := err.(mongo.CommandError); ok {
			return documentsInserted, &libraryErrors.ConnectionError{Db: mongoDb}
		} else if _, ok := err.(mongo.BulkWriteException); ok {
			return documentsInserted, &libraryErrors.AlreadyExistError{Message: "Document with a key already exists"}
		} else {
			return documentsInserted, err
		}
	}

	for _, id := range insertResult.InsertedIDs {
		documentReturned, err := manager.FindOne(map[string]interface{}{"_id": id})
		if err != nil {
			if _, ok := err.(mongo.CommandError); ok {
				return nil, &libraryErrors.ConnectionError{Db: mongoDb}
			} else {
				return nil, err
			}
		}
		documentsInserted = append(documentsInserted, documentReturned)
	}

	return documentsInserted, nil
}

// FindOne is the function to find just one document that matches with the filter
// filter: It is the filter to find the document inside the MongoDB
// It returns the first document matching with the filter and an error
func (manager *MongoManager) FindOne(filter map[string]interface{}) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(manager.timeout)*time.Second)
	defer cancel()

	resultFind := manager.collection.FindOne(ctx, filter)
	if err := resultFind.Err(); err != nil {
		if _, ok := err.(mongo.CommandError); ok {
			return nil, &libraryErrors.ConnectionError{Db: mongoDb}
		} else if reflect.TypeOf(err).String() == "*errors.errorString" {
			return nil, &libraryErrors.NotExistError{Message: "Document not found"}
		} else {
			return nil, err
		}
	}
	var documentReturned bson.M
	var err = resultFind.Decode(&documentReturned)
	return documentReturned, err
}

// Find is the function inside the MongoManager to return a list of documents that match the filter defined
// filter: It is the filter to find documents inside the MongoDB
// It returns a list of documents (it may be empty) and an error
func (manager *MongoManager) FindMany(filter map[string]interface{}) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(manager.timeout)*time.Second)
	defer cancel()

	var results []map[string]interface{}
	cursor, err := manager.collection.Find(ctx, filter)
	if err != nil {
		if _, ok := err.(mongo.CommandError); ok {
			return nil, &libraryErrors.ConnectionError{Db: mongoDb}
		} else {
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
// update: It contains the new changes to apply in the document
// It returns the document updated and an error
func (manager *MongoManager) UpdateOne(filter map[string]interface{}, update interface{}) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(manager.timeout)*time.Second)
	defer cancel()

	m, ok := update.(map[string]interface{})
	if !ok {
		return nil, errors.New("update is not a map[string]interface{}")
	}

	var documentReturned bson.M
	err := manager.collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": m}).Decode(&documentReturned)
	if err != nil {
		if _, ok := err.(mongo.CommandError); ok {
			return nil, &libraryErrors.ConnectionError{Db: mongoDb}
		} else if reflect.TypeOf(err).String() == "*errors.errorString" {
			return nil, &libraryErrors.NotExistError{Message: "Document not found"}
		} else {
			return nil, err
		}
	}

	return manager.FindOne(map[string]interface{}{"_id": documentReturned["_id"]})
}

// UpdateMany is the function for updating multiple documents that match the filter
// filter: It is the filter to find the documents
// update: The new content for the documents found
// It returns the documents updated and an error
func (manager *MongoManager) UpdateMany(filter map[string]interface{}, update interface{}) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(manager.timeout)*time.Second)
	opts := options.Update().SetUpsert(false)
	defer cancel()

	m, ok := update.(map[string]interface{})
	if !ok {
		return nil, errors.New("i is not a map[string]interface{}")
	}
	documentsFilter, err := manager.FindMany(filter)
	if err != nil {
		return nil, err
	}
	if len(documentsFilter) == 0 {
		return nil, &libraryErrors.NotExistError{Message: "Document not found"}
	}

	resultUpdate, err := manager.collection.UpdateMany(ctx, filter, bson.M{"$set": m}, opts)
	if err != nil {
		return nil, err
	}
	if resultUpdate.MatchedCount == 0 {
		return nil, &libraryErrors.NotExistError{Message: "Document not found"}
	}

	var documentsModified []map[string]interface{}
	for _, document := range documentsFilter {
		documentReturned, err := manager.FindOne(map[string]interface{}{"_id": document["_id"]})
		if err != nil {
			if _, ok := err.(mongo.CommandError); ok {
				return nil, &libraryErrors.ConnectionError{Db: mongoDb}
			} else {
				return nil, err
			}
		}
		documentsModified = append(documentsModified, documentReturned)
	}

	return documentsModified, err
}

// DeleteOne is the function inside the MongoManager to delete the first document that matches with the filter
// filter: It is the filter to find the document to delete
// It returns an error in case a document was not deleted
func (manager *MongoManager) DeleteOne(filter map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(manager.timeout)*time.Second)
	defer cancel()

	result, err := manager.collection.DeleteOne(ctx, filter)
	if result.DeletedCount == 0 {
		return &libraryErrors.NotExistError{Message: "Document not found"}
	}
	return err
}

// DeleteMany is the function inside the MongoManager to delete all the documents that match with the filter
// filter: It is the filter to find the documents to delete
// It returns the number of documents deleted and an error
func (manager *MongoManager) DeleteMany(filter map[string]interface{}) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(manager.timeout)*time.Second)
	defer cancel()

	result, err := manager.collection.DeleteMany(ctx, filter)
	return int(result.DeletedCount), err
}
