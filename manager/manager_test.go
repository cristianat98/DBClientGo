package manager

import (
	"os"
	"testing"

	libraryErrors "github.com/cristianat98/dbclientgo/errors"
	"github.com/cristianat98/dbclientgo/mongo"
	"github.com/stretchr/testify/assert"
)

const timeoutTest = 5
const dbTest = "test"

func TestCreateDatabaseManagerSuccess(t *testing.T) {
	mongoUri := os.Getenv("MONGO_URI")
	assert.NotEqual(t, "", mongoUri)

	databaseManager, err := CreateDatabaseManager(mongoUri, dbTest, mongo.MONGO, timeoutTest)
	assert.NotNil(t, databaseManager)
	assert.NoError(t, err)
	databaseManager.DisconnectDb()
}

func TestCreateDatabaseManagerFailedInvalidType(t *testing.T) {
	databaseManager, err := CreateDatabaseManager("test", dbTest, "test", timeoutTest)
	assert.Nil(t, databaseManager)
	var myErr *libraryErrors.InputError
	assert.ErrorAs(t, err, &myErr)
}
