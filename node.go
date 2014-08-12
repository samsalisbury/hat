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
	Node          *Node
	Parent        *ResolvedNode
	ID            string
	Entity        interface{} // Can be collection or single struct
	CollectionIDs []string    // Only set when Node.IsCollection
	Tag           *Tag
	HTTPMethods   map[string]StdHTTPMethod
}

func (root *Node) Locate(path ...string) (*ResolvedNode, error) {
	if resolvedRoot, err := ResolveRoot(root); err != nil {
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
	child, err := n.Resolve(id)
	if err != nil {
		return nil, err
	} else if child == nil {
		return nil, HttpError(404, id, "does not exist")
	}
	return child.Locate(path...)
}

func ResolveRoot(root *Node) (*ResolvedNode, error) {
	entity, _, err := root.Manifest(nil, "")
	if err != nil {
		return nil, err
	}
	return &ResolvedNode{root, nil, "", entity, nil, &Tag{}, nil}, nil
}

func (n *ResolvedNode) Resolve(id string) (*ResolvedNode, error) {
	if n.Node.IsCollection {
		return n.ResolveItem(id)
	} else {
		return n.ResolveMember(id)
	}
}

func (n *Node) Manifest(parentEntity interface{}, id string) (entity interface{}, ids []string, err error) {
	inputs := n.createChildManifestInputs(parentEntity, id)
	var other interface{}
	if n.IsCollection {
		entity, other, err = n.Ops["Page"].Invoke(inputs)
		ids = other.([]string)
	} else {
		entity, _, err = n.Ops["Manifest"].Invoke(inputs)
	}
	return
}

func (n *Node) ManifestStruct(inputs map[IN]boundInput) (entity, _ interface{}, err error) {
	return n.Ops["Manifest"].Invoke(inputs)
}

func (n *Node) ManifestCollection(inputs map[IN]boundInput) (entity, ids interface{}, err error) {
	return n.Ops["Page"].Invoke(inputs)
}

func (n *ResolvedNode) ResolveItem(id string) (*ResolvedNode, error) {
	childEntity, childIDs, err := n.Node.Collection.Node.Manifest(n.Entity, id)
	if err != nil {
		return nil, err
	}
	if childEntity == nil {
		return nil, HttpError(404, "collection", n.ID, "does not have an item with ID", quot(id))
	}
	return n.newResolvedItem(id, childEntity, childIDs), nil
	//return &ResolvedNode{n.Node.Collection.Node, n, id, childEntity, childIDs, n.Node.Collection.Tag, nil}, nil
}

func (n *ResolvedNode) ResolveMember(id string) (*ResolvedNode, error) {
	member, ok := n.Node.Members[id]
	if !ok {
		return nil, HttpError(404, n.ID, "does not have a member called", quot(id))
	}
	childEntity, childIDs, err := member.Node.Manifest(n.Entity, id)
	if err != nil {
		return nil, err
	}
	return n.newResolvedMember(member, id, childEntity, childIDs), nil
}

// Creates both singular and collection resolved nodes that belong to a collection.
func (n *ResolvedNode) newResolvedItem(id string, entity interface{}, collectionIDs []string) *ResolvedNode {
	return &ResolvedNode{n.Node.Collection.Node, n, id, entity, collectionIDs, n.Node.Collection.Tag, nil}
}

// Creates both singular and collection resolved nodes that belong to a named member.
func (n *ResolvedNode) newResolvedMember(member *Member, id string, entity interface{}, collectionIDs []string) *ResolvedNode {
	return &ResolvedNode{member.Node, n, id, entity, collectionIDs, member.Tag, nil}
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
