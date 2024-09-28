package mongo

import (
	"errors"
	"os"
	"testing"

	libraryErrors "github.com/cristianat98/dbclientgo/errors"
	"github.com/stretchr/testify/assert"
)

func initializeDb() (*MongoManager, error) {
	mongoUri := os.Getenv("MONGO_URI")
	if mongoUri == "" {
		return nil, errors.New("MONGO_URI is not set")
	}
	mongoManager, err := CreateMongoManager(mongoUri, "test", "test", 5)
	if err != nil {
		return nil, err
	}

	_, err = mongoManager.DeleteMany(map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	result, err := mongoManager.FindMany(map[string]interface{}{})
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

	insertDocument := map[string]interface{}{
		"test": "test",
	}
	result, err := mongoManager.InsertOne(insertDocument)
	expected := insertDocument
	expected["_id"] = result["_id"]
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestInsertOneFailedIdAlreadyExists(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	insertDocument := map[string]interface{}{
		"test": "test",
	}
	result, err := mongoManager.InsertOne(insertDocument)
	assert.NoError(t, err)
	insertDocument["_id"] = result["_id"]

	result, err = mongoManager.InsertOne(insertDocument)
	assert.Nil(t, result)
	var myErr *libraryErrors.AlreadyExistError
	assert.ErrorAs(t, err, &myErr)
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
	result, err := mongoManager.InsertMany(insertDocuments)
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

	var insertDocuments []map[string]interface{}
	insertDocument1 := map[string]interface{}{
		"document1": "test",
	}
	insertDocuments = append(insertDocuments, insertDocument1)
	insertDocument2 := map[string]interface{}{
		"document2": "test",
	}
	insertDocuments = append(insertDocuments, insertDocument2)
	result, err := mongoManager.InsertMany(insertDocuments)
	assert.NoError(t, err)

	insertDocument1["_id"] = result[0]["_id"]
	result, err = mongoManager.InsertMany(insertDocuments)
	assert.Nil(t, result)
	var myErr *libraryErrors.AlreadyExistError
	assert.ErrorAs(t, err, &myErr)
}

func TestFindOneSuccess(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	insertDocument := map[string]interface{}{
		"test": "test",
	}
	resultInsert, err := mongoManager.InsertOne(insertDocument)
	assert.NoError(t, err)

	resultFind, err := mongoManager.FindOne(insertDocument)
	assert.NoError(t, err)
	assert.Equal(t, resultInsert, resultFind)
}

func TestFindOneFailedNoExist(t *testing.T) {
	mongoManager, err := initializeDb()
	assert.NoError(t, err)

	document := map[string]interface{}{
		"test": "test",
	}

	result, err := mongoManager.FindOne(document)
	assert.Nil(t, result)
	var myErr *libraryErrors.NotExistError
	assert.ErrorAs(t, err, &myErr)
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
	resultInsert, err := mongoManager.InsertMany(insertDocuments)
	assert.NoError(t, err)

	resultFind, err := mongoManager.FindMany(insertDocument)
	assert.NoError(t, err)
	assert.Equal(t, resultInsert, resultFind)
}
