package mongo

import (
	"errors"
	"os"
	"testing"

	libraryErrors "github.com/cristianat98/dbclientgo/errors"
	"github.com/stretchr/testify/assert"
)

const (
	timeoutTest    = 5
	dbTest         = "test"
	collectionTest = "test"
)

func initializeDb() (*MongoManager, error) {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return nil, errors.New("MONGO_URI is not set")
	}
	mongoManager, err := CreateMongoManager(mongoURI, dbTest, timeoutTest)
	if err != nil {
		return nil, err
	}

	_, err = mongoManager.DeleteMany(collectionTest, timeoutTest, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	result, err := mongoManager.FindMany(collectionTest, timeoutTest, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	if len(result) != 0 {
		return nil, errors.New("DB not empty")
	}
	return mongoManager, nil
}

func TestConnectDbSuccess(t *testing.T) {
	mongoURI := os.Getenv("MONGO_URI")
	assert.NotEqual(t, "", mongoURI)

	mongoManager := new(MongoManager)
	err := mongoManager.ConnectDb(mongoURI, dbTest, timeoutTest)
	assert.NoError(t, err)
	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}

func TestConnectDbFailedConnectionError(t *testing.T) {
	mongoManager := new(MongoManager)
	err := mongoManager.ConnectDb("mongodb://test", dbTest, 1)
	var myErr *libraryErrors.ConnectionError
	assert.ErrorAs(t, err, &myErr)
}

func TestInsertOneSuccess(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	insertDocument := map[string]interface{}{
		"test": "test",
	}
	result, err := mongoManager.InsertOne(collectionTest, timeoutTest, insertDocument)
	expected := insertDocument
	expected["_id"] = result["_id"]
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}

func TestInsertOneFailedIdAlreadyExists(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	insertDocument := map[string]interface{}{
		"test": "test",
	}
	result, err := mongoManager.InsertOne(collectionTest, timeoutTest, insertDocument)
	assert.NoError(t, err)
	insertDocument["_id"] = result["_id"]

	result, err = mongoManager.InsertOne(collectionTest, timeoutTest, insertDocument)
	assert.Nil(t, result)
	var myErr *libraryErrors.AlreadyExistError
	assert.ErrorAs(t, err, &myErr)

	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}

func TestInsertOneFailedInvalidTimeout(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	insertDocument := map[string]interface{}{
		"test": "test",
	}
	result, err := mongoManager.InsertOne(collectionTest, 0, insertDocument)
	assert.Nil(t, result)
	var myErr *libraryErrors.InputError
	assert.ErrorAs(t, err, &myErr)

	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}

func TestInsertManySuccess(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	var insertDocuments []map[string]interface{}
	insertDocument1 := map[string]interface{}{
		"document1": "test",
	}
	insertDocuments = append(insertDocuments, insertDocument1)
	insertDocument2 := map[string]interface{}{
		"document2": "test",
	}
	insertDocuments = append(insertDocuments, insertDocument2)
	result, err := mongoManager.InsertMany(collectionTest, timeoutTest, insertDocuments)
	assert.NoError(t, err)
	for index, item := range result {
		expected := insertDocuments[index]
		expected["_id"] = item["_id"]
		assert.Equal(t, expected, item)
	}

	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}

func TestInsertManyFailedIdAlreadyExists(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	var insertDocuments []map[string]interface{}
	insertDocument1 := map[string]interface{}{
		"document1": "test",
	}
	insertDocuments = append(insertDocuments, insertDocument1)
	insertDocument2 := map[string]interface{}{
		"document2": "test",
	}
	insertDocuments = append(insertDocuments, insertDocument2)
	result, err := mongoManager.InsertMany(collectionTest, timeoutTest, insertDocuments)
	assert.NoError(t, err)

	insertDocument2["_id"] = result[1]["_id"]
	result, err = mongoManager.InsertMany(collectionTest, timeoutTest, insertDocuments)
	var expected []map[string]interface{}
	expected = append(expected, insertDocument1)
	expected[0]["_id"] = result[0]["_id"]
	assert.Equal(t, expected, result)
	var myErr *libraryErrors.AlreadyExistError
	assert.ErrorAs(t, err, &myErr)

	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}

func TestInsertManyFailedInvalidTimeout(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	var insertDocuments []map[string]interface{}
	insertDocument1 := map[string]interface{}{
		"document1": "test",
	}
	insertDocuments = append(insertDocuments, insertDocument1)
	insertDocument2 := map[string]interface{}{
		"document2": "test",
	}
	insertDocuments = append(insertDocuments, insertDocument2)
	result, err := mongoManager.InsertMany(collectionTest, 0, insertDocuments)
	assert.Nil(t, result)
	var myErr *libraryErrors.InputError
	assert.ErrorAs(t, err, &myErr)

	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}

func TestFindOneSuccess(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	insertDocument := map[string]interface{}{
		"test": "test",
	}
	resultInsert, err := mongoManager.InsertOne(collectionTest, timeoutTest, insertDocument)
	assert.NoError(t, err)

	resultFind, err := mongoManager.FindOne(collectionTest, timeoutTest, insertDocument)
	assert.NoError(t, err)
	assert.Equal(t, resultInsert, resultFind)

	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}

func TestFindOneFailedNoExist(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	document := map[string]interface{}{
		"test": "test",
	}

	result, err := mongoManager.FindOne(collectionTest, timeoutTest, document)
	assert.Nil(t, result)
	var myErr *libraryErrors.NotExistError
	assert.ErrorAs(t, err, &myErr)

	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}

func TestFindOneFailedInvalidTimeout(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	document := map[string]interface{}{
		"test": "test",
	}
	result, err := mongoManager.FindOne(collectionTest, 0, document)
	assert.Nil(t, result)
	var myErr *libraryErrors.InputError
	assert.ErrorAs(t, err, &myErr)

	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}

func TestFindManySuccess(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	var insertDocuments []map[string]interface{}
	insertDocument := map[string]interface{}{
		"document": "test",
	}
	insertDocuments = append(insertDocuments, insertDocument)
	insertDocuments = append(insertDocuments, insertDocument)
	resultInsert, err := mongoManager.InsertMany(collectionTest, timeoutTest, insertDocuments)
	assert.NoError(t, err)

	resultFind, err := mongoManager.FindMany(collectionTest, timeoutTest, insertDocument)
	assert.NoError(t, err)
	assert.Equal(t, resultInsert, resultFind)

	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}

func TestFindManyFailedInvalidTimeout(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	document := map[string]interface{}{
		"test": "test",
	}
	result, err := mongoManager.FindMany(collectionTest, 0, document)
	assert.Nil(t, result)
	var myErr *libraryErrors.InputError
	assert.ErrorAs(t, err, &myErr)

	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}

func TestUpdateOneSuccess(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	insertDocument := map[string]interface{}{
		"test": "test",
	}
	resultInsert, err := mongoManager.InsertOne(collectionTest, timeoutTest, insertDocument)
	assert.NoError(t, err)

	updateDocument := map[string]interface{}{
		"test": "test2",
	}
	resultUpdate, err := mongoManager.UpdateOne(collectionTest, timeoutTest, insertDocument, updateDocument)
	expected := resultInsert
	expected["test"] = "test2"
	assert.NoError(t, err)
	assert.Equal(t, expected, resultUpdate)

	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}

func TestUpdateOneAddFieldSuccess(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	insertDocument := map[string]interface{}{
		"test": "test",
	}
	resultInsert, err := mongoManager.InsertOne(collectionTest, timeoutTest, insertDocument)
	assert.NoError(t, err)

	updateDocument := map[string]interface{}{
		"test2": "test2",
	}
	resultUpdate, err := mongoManager.UpdateOne(collectionTest, timeoutTest, insertDocument, updateDocument)
	expected := resultInsert
	expected["test2"] = "test2"
	assert.NoError(t, err)
	assert.Equal(t, expected, resultUpdate)

	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}

func TestUpdateOneFailedNoExist(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	result, err := mongoManager.UpdateOne(collectionTest, timeoutTest, map[string]interface{}{"test": "test"}, map[string]interface{}{"test": "test"})
	assert.Nil(t, result)
	var myErr *libraryErrors.NotExistError
	assert.ErrorAs(t, err, &myErr)

	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}

func TestUpdateOneFailedInvalidTimeout(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	result, err := mongoManager.UpdateOne(collectionTest, 0, map[string]interface{}{"test": "test"}, map[string]interface{}{"test": "test"})
	assert.Nil(t, result)
	var myErr *libraryErrors.InputError
	assert.ErrorAs(t, err, &myErr)

	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}

func TestUpdateManySuccess(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	var insertDocuments []map[string]interface{}
	insertDocument := map[string]interface{}{
		"document": "test",
	}
	insertDocuments = append(insertDocuments, insertDocument)
	insertDocuments = append(insertDocuments, insertDocument)
	resultInsert, err := mongoManager.InsertMany(collectionTest, timeoutTest, insertDocuments)
	assert.NoError(t, err)

	updateDocument := map[string]interface{}{
		"document": "test2",
	}
	resultUpdate, err := mongoManager.UpdateMany(collectionTest, timeoutTest, insertDocument, updateDocument)
	assert.NoError(t, err)
	for index, item := range resultInsert {
		expected := item
		expected["document"] = "test2"
		assert.Equal(t, expected, resultUpdate[index])
	}

	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}

func TestUpdateManyFailedInvalidTimeout(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	result, err := mongoManager.UpdateMany(collectionTest, 0, map[string]interface{}{"test": "test"}, map[string]interface{}{"test": "test"})
	assert.Nil(t, result)
	var myErr *libraryErrors.InputError
	assert.ErrorAs(t, err, &myErr)

	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}

func TestDeleteOneSuccess(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	insertDocument := map[string]interface{}{
		"test": "test",
	}
	_, err = mongoManager.InsertOne(collectionTest, timeoutTest, insertDocument)
	assert.NoError(t, err)

	err = mongoManager.DeleteOne(collectionTest, timeoutTest, insertDocument)
	assert.NoError(t, err)

	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}

func TestDeleteOneFailedNoExist(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	err = mongoManager.DeleteOne(collectionTest, timeoutTest, map[string]interface{}{"test": "test"})
	var myErr *libraryErrors.NotExistError
	assert.ErrorAs(t, err, &myErr)

	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}

func TestDeleteOneFailedInvalidTimeout(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	err = mongoManager.DeleteOne(collectionTest, 0, map[string]interface{}{"test": "test"})
	var myErr *libraryErrors.InputError
	assert.ErrorAs(t, err, &myErr)

	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}

func TestDeleteManySuccess(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	var insertDocuments []map[string]interface{}
	insertDocument := map[string]interface{}{
		"document": "test",
	}
	insertDocuments = append(insertDocuments, insertDocument)
	insertDocuments = append(insertDocuments, insertDocument)
	_, err = mongoManager.InsertMany(collectionTest, timeoutTest, insertDocuments)
	assert.NoError(t, err)

	result, err := mongoManager.DeleteMany(collectionTest, timeoutTest, insertDocument)
	assert.NoError(t, err)
	assert.Equal(t, 2, result)

	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}

func TestDeleteManyFailedInvalidTimeout(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	result, err := mongoManager.DeleteMany(collectionTest, 0, map[string]interface{}{"test": "test"})
	assert.Equal(t, 0, result)
	var myErr *libraryErrors.InputError
	assert.ErrorAs(t, err, &myErr)

	err = mongoManager.DisconnectDb()
	assert.NoError(t, err)
}
