package hat

import (
	"reflect"
)

func (n *Node) initMembers() error {
	switch n.EntityType.Kind() {
	case reflect.Struct:
		return n.initStructMembers()
	case reflect.Map, reflect.Slice:
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
	if childNode, err := newNode(n, elementType); err != nil {
		return err
	} else {
		n.Collection = &Member{Node: childNode, Tag: &Tag{}}
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
		if tagData := f.Tag.Get("hat"); tagData != "" {
			tag, err := parseTag(tagData)
			if err != nil {
				return err
			}

			// Embedded items must be pointers. This just makes the code in hat
			// easier to write.
			if tag.Embed && f.Type.Kind() != reflect.Ptr {
				return n.MethodError(f.Name, "is", f.Type, "should be", reflect.PtrTo(f.Type), "because it is tagged embed()")
			}

			if childNode, err := newNode(n, f.Type); err != nil {
				return err
			} else if member, err := newMember(f.Name, childNode, tag); err != nil {
				return err
			} else {
				if childNode.IsCollection {
					// TODO; Don't flip bits like this, it's hard to maintain.
					// make newNode() more intelligent, maybe by passing tag
					// in so it can configure the node correctly.
					childNode.CollectionTag = tag
				}
				n.Members[member.URLName()] = member
			}
		}
	}
	return nil
}
