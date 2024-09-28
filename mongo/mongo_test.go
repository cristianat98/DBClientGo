package mongo

import (
	"errors"
	"os"
	"testing"

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
	_, err = mongoManager.InsertOne(insertDocument)
	assert.NoError(t, err)

	// err = mongoManager.DeleteOne(insertDocument)
	// assert.NoError(t, err)
	// assert.Equal(t, insertedDocument, insertDocument)

	// var id = primitive.NewObjectID()
	// mongoMock := &mocks.DatabaseMock{
	// 	InsertOneFunc: func(document map[string]interface{}) (map[string]interface{}, error) {
	// 		return map[string]interface{}{"_id": id, "name": "test"}, nil
	// 	},
	// 	FindFunc: func(filter map[string]interface{}) ([]map[string]interface{}, error) {
	// 		var sportsBooks []map[string]interface{}
	// 		return sportsBooks, nil
	// 	},
	// }
	// var sportsBookInterface = CreateMockSportsBookInterface(mongoMock)

	// sportsBookCreate := models.SportsBookCreate{
	// 	Name: "test",
	// }
	// sportsBookCreated := models.SportsBook{
	// 	Id:   id,
	// 	Name: "test",
	// }
	// result, err := sportsBookInterface.CreateSportsBook(sportsBookCreate)

	// assert.NoError(t, err)
	// assert.Equal(t, sportsBookCreated, result)
}

// func TestCreateSportsBookFailsSportsBookAlreadyExists(t *testing.T) {
// 	var id = primitive.NewObjectID()
// 	mongoMock := &mocks.DatabaseMock{
// 		FindFunc: func(filter map[string]interface{}) ([]map[string]interface{}, error) {
// 			var sportsBook = map[string]interface{}{"_id": id, "name": "test"}
// 			var sportsBooks []map[string]interface{}
// 			sportsBooks = append(sportsBooks, sportsBook)
// 			return sportsBooks, nil
// 		},
// 	}
// 	var sportsBookInterface = CreateMockSportsBookInterface(mongoMock)

// 	sportsBookCreate := models.SportsBookCreate{
// 		Name: "test",
// 	}

// 	result, err := sportsBookInterface.CreateSportsBook(sportsBookCreate)
// 	var myErr *utils.AlreadyExistError

// 	assert.Error(t, err)
// 	assert.ErrorAs(t, err, &myErr)
// 	assert.Equal(t, models.SportsBook{}, result)
// }

// func TestCreateSportsBookFailsConnectionClosed(t *testing.T) {
// 	mongoMock := &mocks.DatabaseMock{
// 		FindFunc: func(filter map[string]interface{}) ([]map[string]interface{}, error) {
// 			var sportsBooks []map[string]interface{}
// 			err := new(mongo.CommandError)
// 			return sportsBooks, err
// 		},
// 	}
// 	var sportsBookInterface = CreateMockSportsBookInterface(mongoMock)

// 	sportsBookCreate := models.SportsBookCreate{
// 		Name: "test",
// 	}

// 	result, err := sportsBookInterface.CreateSportsBook(sportsBookCreate)
// 	var myErr *mongo.CommandError

// 	assert.Error(t, err)
// 	assert.ErrorAs(t, err, &myErr)
// 	assert.Equal(t, models.SportsBook{}, result)
// }

// func TestGetSportsBooksSuccess(t *testing.T) {
// 	var id = primitive.NewObjectID()
// 	mongoMock := &mocks.DatabaseMock{
// 		FindFunc: func(filter map[string]interface{}) ([]map[string]interface{}, error) {
// 			var sportsBook = map[string]interface{}{"_id": id, "name": "test"}
// 			var sportsBooks []map[string]interface{}
// 			sportsBooks = append(sportsBooks, sportsBook)
// 			return sportsBooks, nil
// 		},
// 	}
// 	var sportsBookInterface = CreateMockSportsBookInterface(mongoMock)

// 	var sportsBooksReturned []models.SportsBook
// 	sportsBookReturned := models.SportsBook{
// 		Id:   id,
// 		Name: "test",
// 	}
// 	sportsBooksReturned = append(sportsBooksReturned, sportsBookReturned)
// 	result, err := sportsBookInterface.GetSportsBook("")

// 	assert.NoError(t, err)
// 	assert.Equal(t, sportsBooksReturned, result)
// }

// func TestGetSportsBookByNameSuccess(t *testing.T) {
// 	var id = primitive.NewObjectID()
// 	mongoMock := &mocks.DatabaseMock{
// 		FindFunc: func(filter map[string]interface{}) ([]map[string]interface{}, error) {
// 			var sportsBook = map[string]interface{}{"_id": id, "name": "test"}
// 			var sportsBooks []map[string]interface{}
// 			sportsBooks = append(sportsBooks, sportsBook)
// 			return sportsBooks, nil
// 		},
// 	}
// 	var sportsBookInterface = CreateMockSportsBookInterface(mongoMock)

// 	var sportsBooksReturned []models.SportsBook
// 	sportsBookReturned := models.SportsBook{
// 		Id:   id,
// 		Name: "test",
// 	}
// 	sportsBooksReturned = append(sportsBooksReturned, sportsBookReturned)
// 	result, err := sportsBookInterface.GetSportsBook("test")

// 	assert.NoError(t, err)
// 	assert.Equal(t, sportsBooksReturned, result)
// }

// func TestGetSportsBookByNameNotExistsSuccess(t *testing.T) {
// 	mongoMock := &mocks.DatabaseMock{
// 		FindFunc: func(filter map[string]interface{}) ([]map[string]interface{}, error) {
// 			var sportsBook []map[string]interface{}
// 			return sportsBook, nil
// 		},
// 	}
// 	var sportsBookInterface = CreateMockSportsBookInterface(mongoMock)
// 	var sportsBookReturned []models.SportsBook = []models.SportsBook{}

// 	result, err := sportsBookInterface.GetSportsBook("test")

// 	assert.NoError(t, err)
// 	assert.Equal(t, sportsBookReturned, result)
// }

// func TestGetSportsBookFailsConnectionClosed(t *testing.T) {
// 	mongoMock := &mocks.DatabaseMock{
// 		FindFunc: func(filter map[string]interface{}) ([]map[string]interface{}, error) {
// 			var sportsBooks []map[string]interface{}
// 			err := new(mongo.CommandError)
// 			return sportsBooks, err
// 		},
// 	}
// 	var sportsBookInterface = CreateMockSportsBookInterface(mongoMock)
// 	var sportsBookReturned []models.SportsBook = []models.SportsBook{}

// 	result, err := sportsBookInterface.GetSportsBook("test")
// 	var myErr *mongo.CommandError

// 	assert.Error(t, err)
// 	assert.ErrorAs(t, err, &myErr)
// 	assert.Equal(t, sportsBookReturned, result)
// }

// func TestUpdateSportsBookSuccess(t *testing.T) {
// 	var id = primitive.NewObjectID()
// 	mongoMock := &mocks.DatabaseMock{
// 		UpdateOneFunc: func(filter map[string]interface{}, document interface{}) (map[string]interface{}, error) {
// 			return map[string]interface{}{"_id": id, "name": "test"}, nil
// 		},
// 		FindFunc: func(filter map[string]interface{}) ([]map[string]interface{}, error) {
// 			var sportsBook = map[string]interface{}{"_id": id, "name": "test"}
// 			var sportsBooks []map[string]interface{}
// 			sportsBooks = append(sportsBooks, sportsBook)
// 			return sportsBooks, nil
// 		},
// 	}
// 	var sportsBookInterface = CreateMockSportsBookInterface(mongoMock)

// 	sportsBookUpdate := models.SportsBookUpdate{
// 		Name: "test",
// 	}
// 	sportsBookUpdated := models.SportsBook{
// 		Id:   id,
// 		Name: "test",
// 	}
// 	result, err := sportsBookInterface.UpdateSportsBook("test", sportsBookUpdate)

// 	assert.NoError(t, err)
// 	assert.Equal(t, sportsBookUpdated, result)
// }

// func TestUpdateSportsBookFailsSportsBookNoExists(t *testing.T) {
// 	mongoMock := &mocks.DatabaseMock{
// 		FindFunc: func(filter map[string]interface{}) ([]map[string]interface{}, error) {
// 			var sportsBooks []map[string]interface{}
// 			return sportsBooks, nil
// 		},
// 	}
// 	var sportsBookInterface = CreateMockSportsBookInterface(mongoMock)

// 	sportsBookUpdate := models.SportsBookUpdate{
// 		Name: "test",
// 	}

// 	result, err := sportsBookInterface.UpdateSportsBook("test", sportsBookUpdate)
// 	var myErr *utils.NotFoundError

// 	assert.Error(t, err)
// 	assert.ErrorAs(t, err, &myErr)
// 	assert.Equal(t, models.SportsBook{}, result)
// }

// func TestUpdateSportsBookFailsConnectionClosed(t *testing.T) {
// 	mongoMock := &mocks.DatabaseMock{
// 		FindFunc: func(filter map[string]interface{}) ([]map[string]interface{}, error) {
// 			var sportsBooks []map[string]interface{}
// 			err := new(mongo.CommandError)
// 			return sportsBooks, err
// 		},
// 	}
// 	var sportsBookInterface = CreateMockSportsBookInterface(mongoMock)

// 	sportsBookUpdate := models.SportsBookUpdate{
// 		Name: "test",
// 	}

// 	result, err := sportsBookInterface.UpdateSportsBook("test", sportsBookUpdate)
// 	var myErr *mongo.CommandError

// 	assert.Error(t, err)
// 	assert.ErrorAs(t, err, &myErr)
// 	assert.Equal(t, models.SportsBook{}, result)
// }

// func TestDeleteSportsBookSuccess(t *testing.T) {
// 	mongoMock := &mocks.DatabaseMock{
// 		DeleteFunc: func(filter map[string]interface{}) (int, error) {
// 			return 1, nil
// 		},
// 	}
// 	var sportsBookInterface = CreateMockSportsBookInterface(mongoMock)

// 	err := sportsBookInterface.DeleteSportsBook("test")

// 	assert.NoError(t, err)
// }

// func TestDeleteSportsBookFailsSportsBookNoExists(t *testing.T) {
// 	mongoMock := &mocks.DatabaseMock{
// 		DeleteFunc: func(filter map[string]interface{}) (int, error) {
// 			return 0, nil
// 		},
// 	}
// 	var sportsBookInterface = CreateMockSportsBookInterface(mongoMock)

// 	err := sportsBookInterface.DeleteSportsBook("test")
// 	var myErr *utils.NotFoundError

// 	assert.Error(t, err)
// 	assert.ErrorAs(t, err, &myErr)
// }

// func TestDeleteSportsBookFailsConnectionClosed(t *testing.T) {
// 	mongoMock := &mocks.DatabaseMock{
// 		DeleteFunc: func(filter map[string]interface{}) (int, error) {
// 			err := new(mongo.CommandError)
// 			return 0, err
// 		},
// 	}
// 	var sportsBookInterface = CreateMockSportsBookInterface(mongoMock)

// 	err := sportsBookInterface.DeleteSportsBook("test")
// 	var myErr *mongo.CommandError

// 	assert.Error(t, err)
// 	assert.ErrorAs(t, err, &myErr)
// }
