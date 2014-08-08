package hat

import (
	"reflect"
)

type Node struct {
	Parent         *Node
	Field          *reflect.StructField
	EntityType     reflect.Type
	EntityPtrType  reflect.Type
	Operations     compiled_ops
	HTTPMethods    map[string]StdHTTPMethod
	Members        map[string]*Member
	Collection     *Member
	CollectionName string
}

type ResolvedNode struct {
	Node   *Node
	Parent *ResolvedNode
	ID     string
	Entity interface{}
}

func (root *Node) Locate(path ...string) (*ResolvedNode, error) {
	if resolvedRoot, err := root.Resolve(nil, ""); err != nil {
		return nil, err
	} else {
		return resolvedRoot.Locate(path...)
	}
}

func (n *ResolvedNode) Locate(path ...string) (*ResolvedNode, error) {
	if len(path) == 0 || (len(path) == 1 && len(path[0]) == 0) {
		return n, nil
	}
	id := path[0]
	path = path[1:]
	if n.Node.Collection != nil {
		if rn, err := n.Node.Collection.Node.Resolve(n, id); err != nil {
			return nil, err
		} else {
			return rn.Locate(path...)
		}
	} else if member, ok := n.Node.Members[id]; ok {
		if rn, err := member.Node.Resolve(n, id); err != nil {
			return nil, err
		} else {
			return rn.Locate(path...)
		}
	} else {
		return nil, HttpError(404, id, "not found.")
	}
}

func (n *Node) Resolve(parentNode *ResolvedNode, id string) (*ResolvedNode, error) {
	var parent interface{}
	if parentNode != nil && parentNode.Entity != nil {
		parent = parentNode.Entity
	}
	println("ABOUT TO GALL GET WITH PARENT === " + Error(parent).Error())
	if _, entity, err := n.HTTPMethods["GET"](n, parent, id, &Payload{}); err != nil {
		return nil, err
	} else {
		println("SETTING RESOLVED NODE ENTITY TO ", Error(entity).Error())
		return &ResolvedNode{n, parentNode, id, entity}, nil
	}
}

func (n *ResolvedNode) ParentEntity() interface{} {
	if n.Parent == nil {
		return nil
	} else {
		return n.Parent.Entity
	}
}

func (n *ResolvedNode) Path() string {
	if n.Parent != nil {
		return n.Parent.Path() + "/" + n.ID
	} else {
		return n.ID
	}
}

func newNode(parent *Node, field *reflect.StructField, entityType reflect.Type) (*Node, error) {
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}
	entityPtrType := reflect.PtrTo(entityType)
	node := &Node{Parent: parent, Field: field, EntityType: entityType, EntityPtrType: entityPtrType}
	if err := node.init(); err != nil {
		return nil, err
	} else {
		return node, nil
	}
}

func (n *Node) init() error {
	if err := n.initOperations(); err != nil {
		return err
	}
	if err := n.initMembers(); err != nil {
		return err
	}
	if err := n.initHttpMethods(); err != nil {
		return err
	}
	return nil
}
