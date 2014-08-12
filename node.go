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
	Parent *ResolvedNode
	Node   *Node
	ID     string
	Tag    *Tag // The member tag for this relationship.

	// Singular-item-specific fields: (nil for collections)
	Entity interface{} // Manifested singular entity.

	// Collection-specific fields: (nil for singular items)
	Collection    interface{} // Manifested collection.
	CollectionIDs []string    // Manifested collection IDs.

	// The HTTP methods for this node; very late bound.
	HTTPMethods map[string]StdHTTPMethod
}

// Creates both singular and collection resolved nodes that belong to a collection.
func (n *ResolvedNode) newResolvedCollection(node *Node, id string, tag *Tag, collection interface{}, ids []string) *ResolvedNode {
	return &ResolvedNode{n, node, id, tag, nil, collection, ids, nil}
}

// Creates both singular and collection resolved nodes that belong to a named member.
func (n *ResolvedNode) newResolvedSingular(node *Node, id string, tag *Tag, entity interface{}) *ResolvedNode {
	return &ResolvedNode{n, node, id, tag, entity, nil, nil, nil}
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
	entity, err := root.ManifestSingular(nil, "")
	if err != nil {
		return nil, err
	}
	return &ResolvedNode{nil, root, "", &Tag{}, entity, nil, nil, nil}, nil
}

// Resolve methods are named like this:
// Resolve(Item|Member)(Singular|Collection)
// We use Item when the parent is a collection
//        Member when the parent is a struct
//        Singular when the child is a struct
//        Collection when the child is a collection.
//
// i.e. Resolve{Child's relationship to parent}{Child's entity mode}

func (n *ResolvedNode) Resolve(id string) (*ResolvedNode, error) {
	if n.Node.IsCollection {
		return n.ResolveItem(id)
	} else {
		return n.ResolveMember(id)
	}
}

func (n *ResolvedNode) ResolveItem(id string) (*ResolvedNode, error) {
	collection := n.Node.Collection
	if collection.Node.IsCollection {
		return n.ResolveItemCollection(collection.Tag, collection.Node, id)
	} else {
		return n.ResolveItemSingular(collection.Tag, collection.Node, id)
	}
}

func (n *ResolvedNode) ResolveMember(id string) (*ResolvedNode, error) {
	member, ok := n.Node.Members[id]
	if !ok {
		return nil, HttpError(404, n.ID, "does not have a member called", quot(id))
	}
	if member.Node.IsCollection {
		return n.ResolveMemberCollection(member.Tag, member.Node, id)
	} else {
		return n.ResolveMemberSingular(member.Tag, member.Node, id)
	}
}

func (n *ResolvedNode) ResolveMemberCollection(tag *Tag, collectionNode *Node, id string) (*ResolvedNode, error) {
	collection, ids, err := collectionNode.ManifestCollection(n.Entity, id)
	if err != nil {
		return nil, err
	}
	return n.newResolvedCollection(collectionNode, id, tag, collection, ids), nil
}

func (n *ResolvedNode) ResolveMemberSingular(tag *Tag, singularNode *Node, id string) (*ResolvedNode, error) {
	entity, err := singularNode.ManifestSingular(n.Entity, id)
	if err != nil {
		return nil, err
	}
	return n.newResolvedSingular(singularNode, id, tag, entity), nil
}

func (n *ResolvedNode) ResolveItemCollection(tag *Tag, collectionNode *Node, id string) (*ResolvedNode, error) {
	collection, ids, err := collectionNode.ManifestCollection(n.Collection, id)
	if err != nil {
		return nil, err
	}
	if collection == nil {
		return nil, HttpError(404, "collection", n.ID, "does not have an item with ID", quot(id))
	}
	return n.newResolvedCollection(collectionNode, id, tag, collection, ids), nil
}

func (n *ResolvedNode) ResolveItemSingular(tag *Tag, singularNode *Node, id string) (*ResolvedNode, error) {
	entity, err := singularNode.ManifestSingular(n.Collection, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, HttpError(404, "collection", n.ID, "does not have an item with ID", quot(id))
	}
	return n.newResolvedSingular(singularNode, id, tag, entity), nil
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

// func (n *ResolvedNode) ParentEntity() interface{} {
// 	if n.Parent == nil {
// 		return nil
// 	} else {
// 		return n.Parent.Entity
// 	}
// }

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
