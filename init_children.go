package hat

import (
	"reflect"
	"strings"
)

func (n *Node) initMembers() error {
	switch n.EntityType.Kind() {
	case reflect.Struct:
		return n.initStructMembers()
	case reflect.Map:
		return n.initMapCollection()
	default:
		panic("Nodes with kind " + n.EntityType.Kind().String() + " are not yet supported.")
	}
}

func (n *Node) initMapCollection() error {
	keyType := n.EntityType.Key()
	elementType := n.EntityType.Elem()
	if keyType.Kind() != reflect.String {
		return Error("Only collections with string keys are currently supported.")
	}
	if childNode, err := newNode(n, nil, elementType); err != nil {
		return err
	} else {
		member := newMember(childNode, "collection()")
		n.Collection = member
		n.CollectionName = elementType.Name()
	}
	return nil
}

func (n *Node) initStructMembers() error {
	n.Members = map[string]*Member{}
	t := n.EntityType
	numFields := t.NumField()
	for i := 0; i < numFields; i++ {
		f := t.Field(i)
		tag := f.Tag.Get("hat")
		switch tag {
		case "embed()":
			name := strings.ToLower(f.Name)
			if childNode, err := newNode(n, &f, f.Type); err != nil {
				return err
			} else {
				member := newMember(childNode, "embed()")
				n.Members[name] = member
			}
		}
	}
	return nil
}
