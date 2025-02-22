# DBClientGo

Repo to host a Go library for creating clients for different DBs (MongoDB, PostgreSQL...). The main idea is to unify the functions, the messages, etc. for all the databases, even some databases work with documents, tables, etc. It unifies only the simple operations: Connect DB, Disconnect DB, Insert data, Get data, Update data and Delete data.

## Structure

The repo contains an Interface called DatabaseInterface (in the database). The different clients must follow this interface and they are considered Managers. Each independent manager should also try to use the same error types and messages to create the lowest possible work when the database is changed in the project.

The current clients integrated are the following, with its respective manager names:
- MongoDB (MongoManager)

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

## Usage

An example of using a Manager is as follows:

```go
// Create the Manager instance (each Manager may have different inputs for the constructor)
// It will try to connect to the DB.
mongoManager, err = CreateMongoManager("dbURI", "dbName", 5)
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
```

## Support

For getting help, please feel free to use the issues on GitHub.

## Roadmap

The following clients to integrate are:
- PostgreSQL

## Contributing

If you are interested on contributing in the code, fork the repository, modify the code as you wish and create a Pull Request to the develop branch of the repository.

## Licence

License used is GPLv3.
