package hat

import (
	"reflect"
)

type Node struct {
	IsCollection   bool
	Parent         *Node
	EntityType     reflect.Type
	EntityPtrType  reflect.Type
	Ops            map[string]*CompiledOp
	Members        map[string]*Member
	Collection     *Member
	CollectionName string
	CollectionTag  *Tag
}

type ResolvedNode struct {
	Node        *Node
	Parent      *ResolvedNode
	ID          string
	Entity      interface{}
	HTTPMethods map[string]StdHTTPMethod
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
	if entity, err := n.innerGET(parentNode, id); err != nil {
		return nil, err
	} else {
		// That last nil is the inputBinder, which only gets set on the target node.
		return &ResolvedNode{n, parentNode, id, entity, nil}, nil
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

func newNode(parent *Node, entityType reflect.Type) (*Node, error) {
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}
	isCollection := false
	if entityType.Kind() == reflect.Slice || entityType.Kind() == reflect.Map {
		isCollection = true
	}
	entityPtrType := reflect.PtrTo(entityType)
	node := &Node{IsCollection: isCollection, Parent: parent, EntityType: entityType, EntityPtrType: entityPtrType}
	if err := node.init(); err != nil {
		return nil, err
	} else {
		return node, nil
	}
}

func (n *Node) init() error {
	if err := n.initOps(); err != nil {
		return err
	}
	if err := n.initMembers(); err != nil {
		return err
	}
	return nil
}
