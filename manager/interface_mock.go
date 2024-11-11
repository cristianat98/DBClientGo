package manager

type DatabaseInterfaceMock struct {
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

func (m *DatabaseInterfaceMock) ConnectDb(dbUri, dbName string, timeout int64) error {
	return m.ConnectDbFunc(dbUri, dbName, timeout)
}

func (m *DatabaseInterfaceMock) DisconnectDb() {
	m.DisconnectDbFunc()
}

func (m *DatabaseInterfaceMock) InsertOne(table string, timeout int64, data map[string]interface{}) (map[string]interface{}, error) {
	return m.InsertOneFunc(table, timeout, data)
}

func (m *DatabaseInterfaceMock) InsertMany(table string, timeout int64, data []map[string]interface{}) ([]map[string]interface{}, error) {
	return m.InsertManyFunc(table, timeout, data)
}

func (m *DatabaseInterfaceMock) FindOne(table string, timeout int64, filter map[string]interface{}) (map[string]interface{}, error) {
	return m.FindOneFunc(table, timeout, filter)
}

func (m *DatabaseInterfaceMock) FindMany(table string, timeout int64, filter map[string]interface{}) ([]map[string]interface{}, error) {
	return m.FindManyFunc(table, timeout, filter)
}

func (m *DatabaseInterfaceMock) UpdateOne(table string, timeout int64, filter map[string]interface{}, newData interface{}) (map[string]interface{}, error) {
	return m.UpdateOneFunc(table, timeout, filter, newData)
}

func (m *DatabaseInterfaceMock) UpdateMany(table string, timeout int64, filter map[string]interface{}, newData interface{}) ([]map[string]interface{}, error) {
	return m.UpdateManyFunc(table, timeout, filter, newData)
}

func (m *DatabaseInterfaceMock) DeleteOne(table string, timeout int64, filter map[string]interface{}) error {
	return m.DeleteOneFunc(table, timeout, filter)
}

func (m *DatabaseInterfaceMock) DeleteMany(table string, timeout int64, filter map[string]interface{}) (int, error) {
	return m.DeleteManyFunc(table, timeout, filter)
}
