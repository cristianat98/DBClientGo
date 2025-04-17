# DBClientGo

Repo to host a Go library for creating clients for different DBs (MongoDB, PostgreSQL...). The main idea is to unify the functions, the messages, etc. for all the databases, even some databases work with documents, tables, etc. It unifies only the simple operations: Connect DB, Disconnect DB, Insert data, Get data, Update data and Delete data.

## Structure

The repo contains an Interface called DatabaseInterface (in the database). The different clients must follow this interface and they are considered Managers. Each independent manager should also try to use the same error types and messages to create the lowest possible work when the database is changed in the project.

The current clients integrated are the following, with its respective manager names:
- MongoDB (Manager)

Each manager has a constructor to create the object, but it is not mandatory to use it, given that it is possible to create the client from "scratch". Also, each manager contains self-tests to make sure the funcionality of each function works as expected.

Different managers contains the different functions:
- Create<DB>Manager: Function to create an instance of the Manager.
- ConnectDB: Function to connect to the DB.
- DisconnectDB: Function to disconnect to the DB.
- InsertOne: Function to insert 1 entry to the DB.
- InsertMany: Function to insert more than 1 entry to the DB.
- FindOne: Function to get data of 1 entry from the DB.
- FindMany: Function to get data of more than 1 entry from the DB.
- UpdateOne: Function to update 1 entry to the DB.
- UpdateMany: Function to update more than 1 entry to the DB.
- DeleteOne: Function to delete 1 entry from the DB.
- DeleteMany: Function to delete more than 1 entry from the DB.
- GetClient: Function to get the native client for using some specific functions of the client. Not specified in the interface because the return is very specific for each DB.

## Usage

An example of using a Manager is as follows:

```go
// Create the Manager instance (each Manager may have different inputs for the constructor)
// It will try to connect to the DB.
mongoManager, err = CreateManager("dbURI", "dbName", 5)
if err != nil {
    // Code when error is raised
}

filter := map[string]interface{}{
    "test": "test",
}
data, err = mongoManager.FindOne("nameCollection", 5, filter)
if err != nil {
    // Code when error is raised
}

newData := map[string]interface{}{
    "test": "test2",
}
data, err = mongoManager.UpdateOne("nameCollection", 5, filter, newData)
if err != nil {
    // Code when error is raised
}

err = mongoManager.DisconnectDb()
data, err = mongoManager.FindOne("nameCollection", 5, filter)
if err != nil {
    // Code when error is raised
}

// It is possible to create a generic Manager
databaseManager, err := := CreateDatabaseManager(mongo.CreateManager("dbURI", "dbName", 1))
if err != nil {
    // Code when error is raised
}
data, err = databaseManager.FindOne("nameCollection", 5, filter)
if err != nil {
    // Code when error is raised
}
```

## Support

For getting help, please feel free to use the issues on GitHub.

## Roadmap

The following clients to integrate are:
- PostgreSQL

## Contributing

If you are interested on contributing in the code, fork the repository, modify the code as you wish and create a Pull Request to the develop branch of the repository.

To run pre-commit, you need to run the following commands:
```sh
pip install pre-commit
go install github.com/lietu/go-pre-commit@v0.1.0
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.0.2
pre-commit run --all-files
```

For running the tests, you need to run the following code:
```sh
# Linux
Mongo_URI=<MONGO-URL> go test -v -cover ./...
# Windows
$env:Mongo_URI = "<MONGO-URL>"
go test -v -cover ./...
```

## Licence

This project is licensed under the terms of the [GNU General Public License v3.0 (GPLv3)](https://www.gnu.org/licenses/gpl-3.0.html).
