package errors

import "fmt"

type ConnectionError struct {
	Db string
}

func (e *ConnectionError) Error() string {
	return fmt.Sprintf("Connection refused from %s", e.Db)
}

type AlreadyExistError struct {
	Message string
}

func (e *AlreadyExistError) Error() string {
	return e.Message
}
