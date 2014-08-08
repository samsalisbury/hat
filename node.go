package hat

import (
	"reflect"
)

type Node struct {
	Parent        *Node
	Field         *reflect.StructField
	EntityType    reflect.Type
	EntityPtrType reflect.Type
}

type resolvedNode struct {
	Node Node
	ID   string
}

func newNode(parent *Node, field *reflect.StructField, entityType reflect.Type) (*Node, error) {
	entityPtrType := reflect.PtrTo(entityType)
	node := &Node{parent, field, entityType, entityPtrType}
	if err := node.init(); error != nil {
		return nil, err
	} else {
		return node
	}
}

func (n *Node) init() error {
	if err := n.AnalyseOperations(); err != nil {
		return err
	}
	if err := n.AnalyseMembers(); err != nil {
		return err
	}
	if err := n.AnalyseCollection(); err != nil {
		return err
	}
	if err := n.AnalyseMethods(); err != nil {
		return err
	}
}
