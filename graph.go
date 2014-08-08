package hat

import (
	"reflect"
)

func Graph(entity interface{}) (*Node, error) {
	t := reflect.TypeOf(entity)
	return graph(t)
}

func graph(t) (*Node, error) {
	if t.Kind() == reflect.Struct {
		return graphStruct(t)
	} else {
		panic("Can only graph structs at the moment.")
	}
}

func graphStruct(t reflect.Type) (*Node, error) {
	if children, err := analyseChildren(t); err != nil {
		return nil, err
	} else if operations, err := analyseOperations(t); err != nil {
		return nil, err
	}
}
