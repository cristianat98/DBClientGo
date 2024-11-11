package mongo

import (
	"errors"
	"os"
	"testing"

	libraryErrors "github.com/cristianat98/dbclientgo/errors"
	"github.com/stretchr/testify/assert"
)

const timeoutTest = 5
const dbTest = "test"
const collectionTest = "test"

func initializeDb() (*MongoManager, error) {
	mongoUri := os.Getenv("MONGO_URI")
	if mongoUri == "" {
		return nil, errors.New("MONGO_URI is not set")
	}
	mongoManager, err := CreateMongoManager(mongoUri, dbTest, timeoutTest)
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

func TestInsertOneSuccess(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)
	defer mongoManager.DisconnectDb()

	insertDocument := map[string]interface{}{
		"test": "test",
	}
	result, err := mongoManager.InsertOne(collectionTest, timeoutTest, insertDocument)
	expected := insertDocument
	expected["_id"] = result["_id"]
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestInsertOneFailedIdAlreadyExists(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)
	defer mongoManager.DisconnectDb()

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
}

func TestInsertOneFailedInvalidTimeout(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)
	defer mongoManager.DisconnectDb()

	insertDocument := map[string]interface{}{
		"test": "test",
	}
	result, err := mongoManager.InsertOne(collectionTest, 0, insertDocument)
	assert.Nil(t, result)
	var myErr *libraryErrors.InputError
	assert.ErrorAs(t, err, &myErr)
}

func TestInsertManySuccess(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)
	defer mongoManager.DisconnectDb()

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
}

func TestInsertManyFailedIdAlreadyExists(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)
	defer mongoManager.DisconnectDb()

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
}

func TestInsertManyFailedInvalidTimeout(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)
	defer mongoManager.DisconnectDb()

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
}

func TestFindOneSuccess(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)
	defer mongoManager.DisconnectDb()

	insertDocument := map[string]interface{}{
		"test": "test",
	}
	resultInsert, err := mongoManager.InsertOne(collectionTest, timeoutTest, insertDocument)
	assert.NoError(t, err)

	resultFind, err := mongoManager.FindOne(collectionTest, timeoutTest, insertDocument)
	assert.NoError(t, err)
	assert.Equal(t, resultInsert, resultFind)
}

func TestFindOneFailedNoExist(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)
	defer mongoManager.DisconnectDb()

	document := map[string]interface{}{
		"test": "test",
	}

	result, err := mongoManager.FindOne(collectionTest, timeoutTest, document)
	assert.Nil(t, result)
	var myErr *libraryErrors.NotExistError
	assert.ErrorAs(t, err, &myErr)
}

func TestFindOneFailedInvalidTimeout(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)
	defer mongoManager.DisconnectDb()

	document := map[string]interface{}{
		"test": "test",
	}
	result, err := mongoManager.FindOne(collectionTest, 0, document)
	assert.Nil(t, result)
	var myErr *libraryErrors.InputError
	assert.ErrorAs(t, err, &myErr)
}

func TestFindManySuccess(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)
	defer mongoManager.DisconnectDb()

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
}

func TestFindManyFailedInvalidTimeout(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)
	defer mongoManager.DisconnectDb()

	document := map[string]interface{}{
		"test": "test",
	}
	result, err := mongoManager.FindMany(collectionTest, 0, document)
	assert.Nil(t, result)
	var myErr *libraryErrors.InputError
	assert.ErrorAs(t, err, &myErr)
}

func TestUpdateOneSuccess(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)
	defer mongoManager.DisconnectDb()

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
}

func TestUpdateOneAddFieldSuccess(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)
	defer mongoManager.DisconnectDb()

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
}

func TestUpdateOneFailedNoExist(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)
	defer mongoManager.DisconnectDb()

	result, err := mongoManager.UpdateOne(collectionTest, timeoutTest, map[string]interface{}{"test": "test"}, map[string]interface{}{"test": "test"})
	assert.Nil(t, result)
	var myErr *libraryErrors.NotExistError
	assert.ErrorAs(t, err, &myErr)
}

func TestUpdateOneFailedInvalidTimeout(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)
	defer mongoManager.DisconnectDb()

	result, err := mongoManager.UpdateOne(collectionTest, 0, map[string]interface{}{"test": "test"}, map[string]interface{}{"test": "test"})
	assert.Nil(t, result)
	var myErr *libraryErrors.InputError
	assert.ErrorAs(t, err, &myErr)
}

func TestUpdateManySuccess(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)
	defer mongoManager.DisconnectDb()

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
}

func TestUpdateManyFailedInvalidTimeout(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)
	defer mongoManager.DisconnectDb()

	result, err := mongoManager.UpdateMany(collectionTest, 0, map[string]interface{}{"test": "test"}, map[string]interface{}{"test": "test"})
	assert.Nil(t, result)
	var myErr *libraryErrors.InputError
	assert.ErrorAs(t, err, &myErr)
}

func TestDeleteOneSuccess(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)
	defer mongoManager.DisconnectDb()

	insertDocument := map[string]interface{}{
		"test": "test",
	}
	_, err = mongoManager.InsertOne(collectionTest, timeoutTest, insertDocument)
	assert.NoError(t, err)

	err = mongoManager.DeleteOne(collectionTest, timeoutTest, insertDocument)
	assert.NoError(t, err)
}

func TestDeleteOneFailedNoExist(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)
	defer mongoManager.DisconnectDb()

	err = mongoManager.DeleteOne(collectionTest, timeoutTest, map[string]interface{}{"test": "test"})
	var myErr *libraryErrors.NotExistError
	assert.ErrorAs(t, err, &myErr)
}

func TestDeleteOneFailedInvalidTimeout(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)
	defer mongoManager.DisconnectDb()

	err = mongoManager.DeleteOne(collectionTest, 0, map[string]interface{}{"test": "test"})
	var myErr *libraryErrors.InputError
	assert.ErrorAs(t, err, &myErr)
}

func TestDeleteManySuccess(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)
	defer mongoManager.DisconnectDb()

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
}

func TestDeleteManyFailedInvalidTimeout(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)
	defer mongoManager.DisconnectDb()

	result, err := mongoManager.DeleteMany(collectionTest, 0, map[string]interface{}{"test": "test"})
	assert.Equal(t, 0, result)
	var myErr *libraryErrors.InputError
	assert.ErrorAs(t, err, &myErr)
}
