package mongo

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	libraryErrors "github.com/cristianat98/dbclientgo/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Manager is the structure to manage the connections and operations to the MongoDB
// client: It is directly the client to the MongoDB
// database: It is the database to connect in MongoDB
type Manager struct {
	client   *mongo.Client
	database *mongo.Database
}

// CreateManager is the constructor for the Manager. If it can not connect to the MongoDB, it will fail
// dbURI: It is the URI to connect to the MongoDB
// dbName: It is the name of the DB inside the MongoDB
// timeout: It is the time to define the timeout inside the Manager
// It returns the Manager instance and an error
func CreateManager(dbURI, dbName string, timeout int64) (*Manager, error) {
	mongoManager := new(Manager)
	if err := mongoManager.ConnectDb(dbURI, dbName, timeout); err != nil {
		return nil, err
	}
	return mongoManager, nil
}

// ConnectDb is the function inside the Manager to connect to the MongoDB
// dbURI: It is the URI to connect to the MongoDB
// dbName: It is the name of the DB inside the MongoDB
// timeout: It is the time to define the timeout inside the Manager
// It returns an error in case there was some error
func (manager *Manager) ConnectDb(dbURI, dbName string, timeout int64) error {
	if timeout < 1 {
		return &libraryErrors.InputError{Message: fmt.Sprintf(timeoutMessage, timeout)}
	}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dbURI).SetServerAPIOptions(serverAPI)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}

	var result bson.M
	if err := client.Database(dbName).RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return &libraryErrors.ConnectionError{Db: mongoDB}
	}
	manager.client = client
	manager.database = client.Database(dbName)
	return nil
}

// DisconnectDb is the function inside the Manager to disconnect from the MongoDB
func (manager *Manager) DisconnectDb() error {
	return manager.client.Disconnect(context.TODO())
}

// InsertOne is the function inside the Manager to insert a document in the collection
// collection: Name of the collection to insert a document
// timeout: It is the time to define the timeout inside the Manager
// document: It is the document to add in the collection
// It returns the new document inserted in the collection and an error
func (manager *Manager) InsertOne(collection string, timeout int64, document map[string]interface{}) (map[string]interface{}, error) {
	if timeout < 1 {
		return nil, &libraryErrors.InputError{Message: fmt.Sprintf(timeoutMessage, timeout)}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	resultInsert, err := manager.database.Collection(collection).InsertOne(ctx, document)
	if err != nil {
		if _, ok := err.(mongo.CommandError); ok {
			return nil, &libraryErrors.ConnectionError{Db: mongoDB}
		} else if _, ok := err.(mongo.WriteException); ok {
			return nil, &libraryErrors.AlreadyExistError{Message: "Document with a key already exists"}
		} else {
			return nil, err
		}
	}

	documentReturned, err := manager.FindOne(collection, timeout, map[string]interface{}{"_id": resultInsert.InsertedID})
	return documentReturned, err
}

// InsertMany is the function inside the Manager to insert many documents in the collection
// collection: Name of the collection to insert many documents
// timeout: It is the time to define the timeout inside the Manager
// documents: It is the list of documents to insert in the collection
// It returns the new documents inserted in the collection and an error
func (manager *Manager) InsertMany(collection string, timeout int64, documents []map[string]interface{}) ([]map[string]interface{}, error) {
	if timeout < 1 {
		return nil, &libraryErrors.InputError{Message: fmt.Sprintf(timeoutMessage, timeout)}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	var documentsParsed []interface{}
	for _, item := range documents {
		documentsParsed = append(documentsParsed, item)
	}

	insertResult, errInsert := manager.database.Collection(collection).InsertMany(ctx, documentsParsed)
	var documentsInserted []map[string]interface{}

	for _, id := range insertResult.InsertedIDs {
		documentReturned, err := manager.FindOne(collection, timeout, map[string]interface{}{"_id": id})
		if err != nil {
			if _, ok := err.(mongo.CommandError); ok {
				return nil, &libraryErrors.ConnectionError{Db: mongoDB}
			}
			return nil, err
		}
		documentsInserted = append(documentsInserted, documentReturned)
	}

	if errInsert != nil {
		if _, ok := errInsert.(mongo.CommandError); ok {
			return documentsInserted, &libraryErrors.ConnectionError{Db: mongoDB}
		} else if _, ok := errInsert.(mongo.BulkWriteException); ok {
			return documentsInserted, &libraryErrors.AlreadyExistError{Message: "Document with a key already exists"}
		}
		return documentsInserted, errInsert

	}

	return documentsInserted, nil
}

// FindOne is the function to find just one document that matches with the filter
// collection: Name of the collection to find a document
// timeout: It is the time to define the timeout inside the Manager
// filter: It is the filter to find the document inside the MongoDB
// It returns the first document matching with the filter and an error
func (manager *Manager) FindOne(collection string, timeout int64, filter map[string]interface{}) (map[string]interface{}, error) {
	if timeout < 1 {
		return nil, &libraryErrors.InputError{Message: fmt.Sprintf(timeoutMessage, timeout)}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	resultFind := manager.database.Collection(collection).FindOne(ctx, filter)
	if err := resultFind.Err(); err != nil {
		if _, ok := err.(mongo.CommandError); ok {
			return nil, &libraryErrors.ConnectionError{Db: mongoDB}
		} else if reflect.TypeOf(err).String() == "*errors.errorString" {
			return nil, &libraryErrors.NotExistError{Message: documentNotFoundMessage}
		} else {
			return nil, err
		}
	}
	var documentReturned bson.M
	err := resultFind.Decode(&documentReturned)
	return documentReturned, err
}

// FindMany is the function inside the Manager to return a list of documents that match the filter defined
// collection: Name of the collection to find many documents
// timeout: It is the time to define the timeout inside the Manager
// filter: It is the filter to find documents inside the MongoDB
// It returns a list of documents (it may be empty) and an error
func (manager *Manager) FindMany(collection string, timeout int64, filter map[string]interface{}) ([]map[string]interface{}, error) {
	if timeout < 1 {
		return nil, &libraryErrors.InputError{Message: fmt.Sprintf(timeoutMessage, timeout)}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	var results []map[string]interface{}
	cursor, err := manager.database.Collection(collection).Find(ctx, filter)
	if err != nil {
		if _, ok := err.(mongo.CommandError); ok {
			return nil, &libraryErrors.ConnectionError{Db: mongoDB}
		}
		return nil, err

	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var result map[string]interface{}

		err = cursor.Decode(&result)
		results = append(results, result)
	}

	return results, err
}

// UpdateOne is the function inside the Manager to update the first document that matches with the filter defined
// collection: Name of the collection to update a document
// timeout: It is the time to define the timeout inside the Manager
// filter: It is the filter to find documents inside the MongoDB to update
// update: It contains the new changes to apply in the document
// It returns the document updated and an error
func (manager *Manager) UpdateOne(collection string, timeout int64, filter map[string]interface{}, update interface{}) (map[string]interface{}, error) {
	if timeout < 1 {
		return nil, &libraryErrors.InputError{Message: fmt.Sprintf(timeoutMessage, timeout)}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	m, ok := update.(map[string]interface{})
	if !ok {
		return nil, errors.New("update is not a map[string]interface{}")
	}

	var documentReturned bson.M
	err := manager.database.Collection(collection).FindOneAndUpdate(ctx, filter, bson.M{"$set": m}).Decode(&documentReturned)
	if err != nil {
		if _, ok := err.(mongo.CommandError); ok {
			return nil, &libraryErrors.ConnectionError{Db: mongoDB}
		} else if reflect.TypeOf(err).String() == "*errors.errorString" {
			return nil, &libraryErrors.NotExistError{Message: documentNotFoundMessage}
		}
		return nil, err
	}

	return manager.FindOne(collection, timeout, map[string]interface{}{"_id": documentReturned["_id"]})
}

// UpdateMany is the function for updating multiple documents that match the filter
// collection: Name of the collection to update many documents
// timeout: It is the time to define the timeout inside the Manager
// filter: It is the filter to find the documents
// update: The new content for the documents found
// It returns the documents updated and an error
func (manager *Manager) UpdateMany(collection string, timeout int64, filter map[string]interface{}, update interface{}) ([]map[string]interface{}, error) {
	if timeout < 1 {
		return nil, &libraryErrors.InputError{Message: fmt.Sprintf(timeoutMessage, timeout)}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	opts := options.Update().SetUpsert(false)
	defer cancel()

	m, ok := update.(map[string]interface{})
	if !ok {
		return nil, errors.New("i is not a map[string]interface{}")
	}
	documentsFilter, err := manager.FindMany(collection, timeout, filter)
	if err != nil {
		return nil, err
	}
	if len(documentsFilter) == 0 {
		return nil, &libraryErrors.NotExistError{Message: documentNotFoundMessage}
	}

	resultUpdate, err := manager.database.Collection(collection).UpdateMany(ctx, filter, bson.M{"$set": m}, opts)
	if err != nil {
		return nil, err
	}
	if resultUpdate.MatchedCount == 0 {
		return nil, &libraryErrors.NotExistError{Message: documentNotFoundMessage}
	}

	var documentsModified []map[string]interface{}
	for _, document := range documentsFilter {
		documentReturned, err := manager.FindOne(collection, timeout, map[string]interface{}{"_id": document["_id"]})
		if err != nil {
			if _, ok := err.(mongo.CommandError); ok {
				return nil, &libraryErrors.ConnectionError{Db: mongoDB}
			}
			return nil, err
		}
		documentsModified = append(documentsModified, documentReturned)
	}

	return documentsModified, err
}

// DeleteOne is the function inside the Manager to delete the first document that matches with the filter
// collection: Name of the collection to delete a document
// timeout: It is the time to define the timeout inside the Manager
// filter: It is the filter to find the document to delete
// It returns an error in case a document was not deleted
func (manager *Manager) DeleteOne(collection string, timeout int64, filter map[string]interface{}) error {
	if timeout < 1 {
		return &libraryErrors.InputError{Message: fmt.Sprintf(timeoutMessage, timeout)}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	result, err := manager.database.Collection(collection).DeleteOne(ctx, filter)
	if result.DeletedCount == 0 {
		return &libraryErrors.NotExistError{Message: documentNotFoundMessage}
	}
	return err
}

// DeleteMany is the function inside the Manager to delete all the documents that match with the filter
// collection: Name of the collection to delete many documents
// timeout: It is the time to define the timeout inside the Manager
// filter: It is the filter to find the documents to delete
// It returns the number of documents deleted and an error
func (manager *Manager) DeleteMany(collection string, timeout int64, filter map[string]interface{}) (int, error) {
	if timeout < 1 {
		return 0, &libraryErrors.InputError{Message: fmt.Sprintf(timeoutMessage, timeout)}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	result, err := manager.database.Collection(collection).DeleteMany(ctx, filter)
	return int(result.DeletedCount), err
}
