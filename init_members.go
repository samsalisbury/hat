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
		n.Collection = &Member{Node: childNode, Tag: &MemberTag{}}
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
		if tag := f.Tag.Get("hat"); tag != "" {
			name := strings.ToLower(f.Name)
			if childNode, err := newNode(n, &f, f.Type); err != nil {
				return err
			} else if member, err := newMember(childNode, tag); err != nil {
				return err
			} else {
				n.Members[name] = member
			}
		}
	}
	return nil
}
