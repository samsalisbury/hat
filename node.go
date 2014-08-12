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

type ResolvedNode interface {
	Locate(path ...string) (ResolvedNode, error)
	Resolve(id string) (ResolvedNode, error)
	ID() string
	Path() string
	Parent() ResolvedNode
	Entity() interface{}
	Links() ([]Link, error)
	Resource() (*Resource, error)
	EmbeddedResource(*Tag) (*Resource, error)
	// Can these 2 be removed??
	EmbeddedCollectionItems() ([]*Resource, error)
	EmbeddedMembers() (map[string]*Resource, error)
}

type ResolvedNodeBase struct {
	parent ResolvedNode
	Node   *Node
	id     string
	Tag    *Tag // The member tag for this relationship.
	// The HTTP methods for this node; very late bound.
	HTTPMethods map[string]StdHTTPMethod
}

func newResolvedNodeBase(parent ResolvedNode, node *Node, id string, tag *Tag) ResolvedNodeBase {
	return ResolvedNodeBase{parent, node, id, tag, nil}
}

func (n *ResolvedNodeBase) ID() string {
	return n.id
}

func (n *ResolvedNodeBase) Path() string {
	if n.parent != nil {
		return n.parent.Path() + "/" + n.ID()
	} else {
		return n.ID()
	}
}

func (n *ResolvedNodeBase) Parent() ResolvedNode {
	return n.parent
}

func LocateFromRoot(root *Node, path ...string) (ResolvedNode, error) {
	if resolvedRoot, err := ResolveRoot(root); err != nil {
		return nil, err
	} else {
		return resolvedRoot.Locate(path...)
	}
}

func ResolveRoot(root *Node) (ResolvedNode, error) {
	entity, err := root.ManifestSingular(nil, "")
	if err != nil {
		return nil, err
	}
	return newResolvedSingular(nil, root, "", nil, entity), nil
}

func (n *Node) ManifestSingular(parentEntity interface{}, id string) (interface{}, error) {
	inputs := n.createChildManifestInputs(parentEntity, id)
	entity, _, err := n.Ops["Manifest"].Invoke(inputs)
	return entity, err
}

func (n *Node) ManifestCollection(parentCollection interface{}, id string) (collection interface{}, ids []string, err error) {
	inputs := n.createChildManifestInputs(parentCollection, id)
	entity, other, err := n.Ops["Page"].Invoke(inputs)
	if err != nil {
		return nil, nil, err
	}
	ids = other.([]string)
	return entity, ids, err
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
