//go:build test
// +build test

package mongo

type MongoManagerMock struct {
	InsertOneFunc  func(document map[string]interface{}) (map[string]interface{}, error)
	InsertManyFunc func(documents []map[string]interface{}) ([]map[string]interface{}, error)
	FindOneFunc    func(filter map[string]interface{}) (map[string]interface{}, error)
	FindManyFunc   func(filter map[string]interface{}) ([]map[string]interface{}, error)
	UpdateOneFunc  func(filter map[string]interface{}, update interface{}) (map[string]interface{}, error)
	UpdateManyFunc func(filter map[string]interface{}, update interface{}) ([]map[string]interface{}, error)
	DeleteOneFunc  func(filter map[string]interface{}) error
	DeleteManyFunc func(filter map[string]interface{}) (int, error)
}

func (m *MongoCollectionMock) InsertOne(document map[string]interface{}) (map[string]interface{}, error) {
	return m.InsertOneFunc(document)
}

func (m *MongoCollectionMock) InsertMany(documents []map[string]interface{}) ([]map[string]interface{}, error) {
	return m.InsertManyFunc(documents)
}

func (m *MongoCollectionMock) FindOne(filter map[string]interface{}) (map[string]interface{}, error) {
	return m.FindOneFunc(filter)
}

func (m *MongoCollectionMock) FindMany(filter map[string]interface{}) ([]map[string]interface{}, error) {
	return m.FindManyFunc(filter)
}

func (m *MongoCollectionMock) UpdateOne(filter map[string]interface{}, update interface{}) (map[string]interface{}, error) {
	return m.UpdateOneFunc(filter, update)
}

func (m *MongoCollectionMock) UpdateMany(filter map[string]interface{}, update interface{}) ([]map[string]interface{}, error) {
	return m.UpdateManyFunc(filter, update)
}

func (m *MongoCollectionMock) DeleteOne(filter map[string]interface{}) error {
	return m.DeleteOneFunc(filter)
}

func (m *MongoCollectionMock) DeleteMany(filter map[string]interface{}) (int, error) {
	return m.DeleteManyFunc(filter)
}
