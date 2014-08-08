package hat

import (
	"fmt"
	"strings"
)

type hatError struct {
	Message string
}

func (h hatError) Error() string {
	return h.Message
}

func Error(args ...interface{}) error {
	message := []string{}
	for _, a := range args {
		message = append(message, fmt.Sprint(a))
	}
	return hatError{strings.Join(message, " ")}
}
