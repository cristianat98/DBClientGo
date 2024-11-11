package manager

import (
	"fmt"

	libraryErrors "github.com/cristianat98/dbclientgo/errors"
	"github.com/cristianat98/dbclientgo/mongo"
)

// Interface that all the DBs Manager must follow
type DatabaseInterface interface {
	ConnectDb(dbUri, dbName string, timeout int64) error
	DisconnectDb()
	InsertOne(table string, timeout int64, data map[string]interface{}) (map[string]interface{}, error)
	InsertMany(table string, timeout int64, data []map[string]interface{}) ([]map[string]interface{}, error)
	FindOne(table string, timeout int64, filter map[string]interface{}) (map[string]interface{}, error)
	FindMany(table string, timeout int64, filter map[string]interface{}) ([]map[string]interface{}, error)
	UpdateOne(table string, timeout int64, filter map[string]interface{}, newData interface{}) (map[string]interface{}, error)
	UpdateMany(table string, timeout int64, filter map[string]interface{}, newData interface{}) ([]map[string]interface{}, error)
	DeleteOne(table string, timeout int64, filter map[string]interface{}) error
	DeleteMany(table string, timeout int64, filter map[string]interface{}) (int, error)
}

// Structure to define the DatabaseManager that will integrate all the DB Managers
// DBClient: Client Database for the different DBs
type DatabaseManager struct {
	dbClient DatabaseInterface
}

// CreateDatabaseManager is the constructor for the DatabaseManager. If it can not connect to the DB, it will fail
// dbUri: It is the URI to connect to the DB
// dbName: It is the name of the DB
// dbType: It is the type of the Database
// It returns the DatabseManager instance and an error
func CreateDatabaseManager(dbUri, dbName, dbType string, timeout int64) (*DatabaseManager, error) {
	var databaseManager = new(DatabaseManager)

	if dbType == mongo.MONGO {
		dbClient, err := mongo.CreateMongoManager(dbUri, dbName, timeout)
		if err != nil {
			return nil, err
		}
		databaseManager.dbClient = dbClient
	} else {
		return nil, &libraryErrors.NotExistError{Message: fmt.Sprintf("DB of type %s does not exist", dbType)}
	}

	return databaseManager, nil
}

// ConnectDb is the function inside the DatabaseManager to connect to the DB
// dbUri: It is the URI to connect to the DB
// dbName: It is the name of the DB
// timeout: It is the time to define the timeout inside the DatabaseManager
// It returns an error in case there was some error
func (manager *DatabaseManager) ConnectDb(dbUri, dbName string, timeout int64) error {
	return manager.dbClient.ConnectDb(dbUri, dbName, timeout)
}

// DisconnectDb is the function inside the DatabaseManager to disconnect from the DB
func (manager *DatabaseManager) DisconnectDb() {
	manager.dbClient.DisconnectDb()
}

// InsertOne is the function inside the DatabaseManager to insert data in the DB
// table: It is the table/collection to add the new data
// timeout: It is the time to define the timeout inside the DatabaseManager
// data: It is the data to add in the collection
// It returns the new data (row/document) inserted in the DB and an error
func (manager *DatabaseManager) InsertOne(table string, timeout int64, data map[string]interface{}) (map[string]interface{}, error) {
	return manager.dbClient.InsertOne(table, timeout, data)
}

// InsertMany is the function inside the DatabaseManager to insert multiple data in the DB
// table: It is the table/collection to add the new data
// timeout: It is the time to define the timeout inside the DatabaseManager
// data: It is the list of data to insert in the DB
// It returns the new data (rows/documents) inserted in the DB and an error
func (manager *DatabaseManager) InsertMany(table string, timeout int64, data []map[string]interface{}) ([]map[string]interface{}, error) {
	return manager.dbClient.InsertMany(table, timeout, data)
}

// FindOne is the function to find just one document that matches with the filter
// table: It is the table/collection to find data
// timeout: It is the time to define the timeout inside the DatabaseManager
// filter: It is the filter to find the data inside the DB
// It returns the first row/document matching with the filter and an error
func (manager *DatabaseManager) FindOne(table string, timeout int64, filter map[string]interface{}) (map[string]interface{}, error) {
	return manager.dbClient.FindOne(table, timeout, filter)
}

// FindMany is the function inside the DatabaseManager to return a list of documents that match the filter defined
// table: It is the table/collection to find data
// timeout: It is the time to define the timeout inside the DatabaseManager
// filter: It is the filter to find data inside the DB
// It returns a list of rows/documents (it may be empty) and an error
func (manager *DatabaseManager) FindMany(table string, timeout int64, filter map[string]interface{}) ([]map[string]interface{}, error) {
	return manager.dbClient.FindMany(table, timeout, filter)
}

// UpdateOne is the function inside the DatabaseManager to update the first document that matches with the filter defined
// table: It is the table/collection to update the new data
// timeout: It is the time to define the timeout inside the DatabaseManager
// filter: It is the filter to find data inside the DB to update
// newData: It contains the new changes to apply in the row/document
// It returns the row/document updated and an error
func (manager *DatabaseManager) UpdateOne(table string, timeout int64, filter map[string]interface{}, newData interface{}) (map[string]interface{}, error) {
	return manager.dbClient.UpdateOne(table, timeout, filter, newData)
}

// UpdateMany is the function for updating data that match the filter
// table: It is the table/collection to update the new data
// timeout: It is the time to define the timeout inside the DatabaseManager
// filter: It is the filter to find the rows/documents
// newData:  It contains the new changes to apply in the rows/documents
// It returns the rows/documents updated and an error
func (manager *DatabaseManager) UpdateMany(table string, timeout int64, filter map[string]interface{}, newData interface{}) ([]map[string]interface{}, error) {
	return manager.dbClient.UpdateMany(table, timeout, filter, newData)
}

// DeleteOne is the function inside the DatabaseManager to delete the first row/collection that matches with the filter
// table: It is the table/collection to delete data
// timeout: It is the time to define the timeout inside the DatabaseManager
// filter: It is the filter to find the row/document to delete
// It returns an error in case a row/document was not deleted
func (manager *DatabaseManager) DeleteOne(table string, timeout int64, filter map[string]interface{}) error {
	return manager.dbClient.DeleteOne(table, timeout, filter)
}

// DeleteMany is the function inside the DatabaseManager to delete all the data that match with the filter
// table: It is the table/collection to delete data
// timeout: It is the time to define the timeout inside the DatabaseManager
// filter: It is the filter to find the rows/documents to delete
// It returns the number of rows/documents deleted and an error
func (manager *DatabaseManager) DeleteMany(table string, timeout int64, filter map[string]interface{}) (int, error) {
	return manager.dbClient.DeleteMany(table, timeout, filter)
}
