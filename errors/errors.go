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

type NotExistError struct {
	Message string
}

func (e *NotExistError) Error() string {
	return e.Message
}

type InputError struct {
	Message string
}

func (e *InputError) Error() string {
	return e.Message
}
