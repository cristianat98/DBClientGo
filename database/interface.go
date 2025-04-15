package database

// Interface that all the DBs Manager must follow
type DatabaseInterface[T any] interface {
	ConnectDb(dbURI, dbName string, timeout int64) error
	DisconnectDb() error
	InsertOne(table string, timeout int64, data map[string]interface{}) (map[string]interface{}, error)
	InsertMany(table string, timeout int64, data []map[string]interface{}) ([]map[string]interface{}, error)
	FindOne(table string, timeout int64, filter map[string]interface{}) (map[string]interface{}, error)
	FindMany(table string, timeout int64, filter map[string]interface{}) ([]map[string]interface{}, error)
	UpdateOne(table string, timeout int64, filter map[string]interface{}, newData interface{}) (map[string]interface{}, error)
	UpdateMany(table string, timeout int64, filter map[string]interface{}, newData interface{}) ([]map[string]interface{}, error)
	DeleteOne(table string, timeout int64, filter map[string]interface{}) error
	DeleteMany(table string, timeout int64, filter map[string]interface{}) (int, error)
	GetClient() T
}

func CreateDatabaseManager[T any](manager DatabaseInterface[T], err error) (DatabaseInterface[T], error) {
	return manager, err
}
