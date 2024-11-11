package manager

type DatabaseManagerMock struct {
	ConnectDbFunc    func(dbUri, dbName string, timeout int64) error
	DisconnectDbFunc func()
	InsertOneFunc    func(table string, timeout int64, data map[string]interface{}) (map[string]interface{}, error)
	InsertManyFunc   func(table string, timeout int64, data []map[string]interface{}) ([]map[string]interface{}, error)
	FindOneFunc      func(table string, timeout int64, filter map[string]interface{}) (map[string]interface{}, error)
	FindManyFunc     func(table string, timeout int64, filter map[string]interface{}) ([]map[string]interface{}, error)
	UpdateOneFunc    func(table string, timeout int64, filter map[string]interface{}, newData interface{}) (map[string]interface{}, error)
	UpdateManyFunc   func(table string, timeout int64, filter map[string]interface{}, newData interface{}) ([]map[string]interface{}, error)
	DeleteOneFunc    func(table string, timeout int64, filter map[string]interface{}) error
	DeleteManyFunc   func(table string, timeout int64, filter map[string]interface{}) (int, error)
}

func (m *DatabaseManagerMock) ConnectDb(dbUri, dbName string, timeout int64) error {
	return m.ConnectDbFunc(dbUri, dbName, timeout)
}

func (m *DatabaseManagerMock) DisconnectDb() {
	m.DisconnectDbFunc()
}

func (m *DatabaseManagerMock) InsertOne(table string, timeout int64, data map[string]interface{}) (map[string]interface{}, error) {
	return m.InsertOneFunc(table, timeout, data)
}

func (m *DatabaseManagerMock) InsertMany(table string, timeout int64, data []map[string]interface{}) ([]map[string]interface{}, error) {
	return m.InsertManyFunc(table, timeout, data)
}

func (m *DatabaseManagerMock) FindOne(table string, timeout int64, filter map[string]interface{}) (map[string]interface{}, error) {
	return m.FindOneFunc(table, timeout, filter)
}

func (m *DatabaseManagerMock) FindMany(table string, timeout int64, filter map[string]interface{}) ([]map[string]interface{}, error) {
	return m.FindManyFunc(table, timeout, filter)
}

func (m *DatabaseManagerMock) UpdateOne(table string, timeout int64, filter map[string]interface{}, newData interface{}) (map[string]interface{}, error) {
	return m.UpdateOneFunc(table, timeout, filter, newData)
}

func (m *DatabaseManagerMock) UpdateMany(table string, timeout int64, filter map[string]interface{}, newData interface{}) ([]map[string]interface{}, error) {
	return m.UpdateManyFunc(table, timeout, filter, newData)
}

func (m *DatabaseManagerMock) DeleteOne(table string, timeout int64, filter map[string]interface{}) error {
	return m.DeleteOneFunc(table, timeout, filter)
}

func (m *DatabaseManagerMock) DeleteMany(table string, timeout int64, filter map[string]interface{}) (int, error) {
	return m.DeleteManyFunc(table, timeout, filter)
}
