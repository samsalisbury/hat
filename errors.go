package hat

import (
	"fmt"
	"strings"
)

type hatError struct {
	Message string
}

type HTTPError interface {
	error
	StatusCode() int
	Err() error
}

type httpError struct {
	statusCode int
	hatError
}

func (he httpError) StatusCode() int {
	return he.statusCode
}

func (he httpError) Err() error {
	return he.hatError
}

func (h hatError) Error() string {
	return h.Message
}

func Error(args ...interface{}) hatError {
	message := []string{}
	for _, a := range args {
		message = append(message, fmt.Sprint(a))
	}
	return hatError{strings.Join(message, " ")}
}

func HttpError(statusCode int, args ...interface{}) HTTPError {
	return httpError{statusCode, Error(args...)}
}

func (n *Node) Error(args ...interface{}) hatError {
	args = append([]interface{}{n.EntityType.Name()}, args...)
	return Error(args...)
}

func (n *Node) MethodError(name string, args ...interface{}) hatError {
	args = append([]interface{}{n.EntityType.Name() + "." + name}, args...)
	return Error(args...)
}

func debug(args ...interface{}) {
	println(Error(args...).Error())
}
