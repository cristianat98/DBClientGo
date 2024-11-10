package manager

import (
	"fmt"

	libraryErrors "github.com/cristianat98/dbclientgo/errors"
	"github.com/cristianat98/dbclientgo/mongo"
)

const MONGO = "MONGO"

// Interface that all the DBs Manager must follow
type DatabaseInterface interface {
	ConnectDb(dbUri, dbName, collection string, timeout int64) error
	DisconnectDb()
	InsertOne(data map[string]interface{}) (map[string]interface{}, error)
	InsertMany(data []map[string]interface{}) ([]map[string]interface{}, error)
	FindOne(filter map[string]interface{}) (map[string]interface{}, error)
	FindMany(filter map[string]interface{}) ([]map[string]interface{}, error)
	UpdateOne(filter map[string]interface{}, newData interface{}) (map[string]interface{}, error)
	UpdateMany(filter map[string]interface{}, newData interface{}) ([]map[string]interface{}, error)
	DeleteOne(filter map[string]interface{}) error
	DeleteMany(filter map[string]interface{}) (int, error)
}

// Structure to define some high-level options to the DB Manager
// collection: Name of the collection when it is needed (MONGO)
// timeout: Timeout in seconds before consider a ConnectionError
// rollback: Bool to define if a rollback will be applied when multiple actions must be done
type Opts struct {
	collection string
	timeout    int64
}

// Structure to define the DatabaseManager that will integrate all the DB Managers
// DBClient: Client Database for the different DBs
// Opts: Configuration for the DatabaseManager
type DatabaseManager struct {
	DbClient DatabaseInterface
	Opts     Opts
}

// CreateDatabaseManager is the constructor for the DatabaseManager. If it can not connect to the DB, it will fail
// dbUri: It is the URI to connect to the DB
// dbName: It is the name of the DB
// dbType: It is the type of the Database
// opts: It defines some extra configuration if needed
// It returns the DatabseManager instance and an error
func CreateDatabaseManager(dbUri, dbName, dbType string, opts Opts) (*DatabaseManager, error) {
	var databaseManager = new(DatabaseManager)
	databaseManager.Opts = opts

	if dbType == MONGO {
		dbClient, err := mongo.CreateMongoManager(dbUri, dbName, opts.collection, opts.timeout)
		if err != nil {
			return nil, err
		}
		databaseManager.DbClient = dbClient
	} else {
		return nil, &libraryErrors.NotExistError{Message: fmt.Sprintf("DB of type %s does not exist", dbType)}
	}

	return databaseManager, nil
}

// ConnectDb is the function inside the DatabaseManager to connect to the DB
// dbUri: It is the URI to connect to the DB
// dbName: It is the name of the DB
// collection: It is the name of the collection inside the DB
// timeout: It is the time to define the timeout inside the DatabaseManager
// It returns an error in case there was some error
func (manager *DatabaseManager) ConnectDb(dbUri, dbName, collection string, timeout int64) error {
	return manager.DbClient.ConnectDb(dbUri, dbName, collection, timeout)
}

// DisconnectDb is the function inside the DatabaseManager to disconnect from the DB
func (manager *DatabaseManager) DisconnectDb() {
	manager.DbClient.DisconnectDb()
}

// InsertOne is the function inside the DatabaseManager to insert data in the DB
// document: It is the document to add in the collection
// It returns the new document inserted in the DB and an error
func (manager *DatabaseManager) InsertOne(document map[string]interface{}) (map[string]interface{}, error) {
	return manager.DbClient.InsertOne(document)
}

// InsertMany is the function inside the DatabaseManager to insert multiple data in the DB
// documents: It is the list of documents to insert in the DB
// It returns the new data inserted in the DB and an error
func (manager *DatabaseManager) InsertMany(documents []map[string]interface{}) ([]map[string]interface{}, error) {
	return manager.DbClient.InsertMany(documents)
}

// FindOne is the function to find just one document that matches with the filter
// filter: It is the filter to find the document inside the DB
// It returns the first document matching with the filter and an error
func (manager *DatabaseManager) FindOne(filter map[string]interface{}) (map[string]interface{}, error) {
	return manager.DbClient.FindOne(filter)
}

// FindMany is the function inside the DatabaseManager to return a list of documents that match the filter defined
// filter: It is the filter to find documents inside the DB
// It returns a list of documents (it may be empty) and an error
func (manager *DatabaseManager) FindMany(filter map[string]interface{}) ([]map[string]interface{}, error) {
	return manager.DbClient.FindMany(filter)
}

// UpdateOne is the function inside the DatabaseManager to update the first document that matches with the filter defined
// filter: It is the filter to find documents inside the DB to update
// update: It contains the new changes to apply in the document
// It returns the document updated and an error
func (manager *DatabaseManager) UpdateOne(filter map[string]interface{}, update interface{}) (map[string]interface{}, error) {
	return manager.DbClient.UpdateOne(filter, update)
}

// UpdateMany is the function for updating multiple documents that match the filter
// filter: It is the filter to find the documents
// update: The new content for the documents found
// It returns the documents updated and an error
func (manager *DatabaseManager) UpdateMany(filter map[string]interface{}, update interface{}) ([]map[string]interface{}, error) {
	return manager.DbClient.UpdateMany(filter, update)
}

// DeleteOne is the function inside the DatabaseManager to delete the first document that matches with the filter
// filter: It is the filter to find the document to delete
// It returns an error in case a document was not deleted
func (manager *DatabaseManager) DeleteOne(filter map[string]interface{}) error {
	return manager.DbClient.DeleteOne(filter)
}

// DeleteMany is the function inside the DatabaseManager to delete all the documents that match with the filter
// filter: It is the filter to find the documents to delete
// It returns the number of documents deleted and an error
func (manager *DatabaseManager) DeleteMany(filter map[string]interface{}) (int, error) {
	return manager.DbClient.DeleteMany(filter)
}
