package hat

import "reflect"

type StdHTTPMethod func(inputs map[IN]boundInput) (statusCode int, resource interface{}, err error)

func makeHTTPMethods(n *ResolvedNode) map[string]StdHTTPMethod {
	return map[string]StdHTTPMethod{
		"GET": makeGET(n),
	}
}

func makeGET(n *ResolvedNode) StdHTTPMethod {
	if n.Node.IsCollection {
		return makeGETCollection(n)
	} else {
		return makeGETMember(n)
	}
}

func makeGETCollection(n *ResolvedNode) StdHTTPMethod {
	return func(inputs map[IN]boundInput) (statusCode int, entity interface{}, err error) {
		entity, _, err = n.Node.Ops["Page"].Invoke(inputs)
		if err != nil {
			return 0, nil, err
		} else if entity == nil {
			return 0, nil, Error("page must set the receiver to a non-nil collection")
		} else {
			if entity, err := n.Node.innerGET(n, n.ID); err != nil {
				return 0, nil, err
			} else {
				return 200, entity, nil
			}
		}
	}
}

func makeGETMember(n *ResolvedNode) StdHTTPMethod {
	return func(inputs map[IN]boundInput) (statusCode int, entity interface{}, err error) {
		if entity, _, err := n.Node.Ops["Manifest"].Invoke(inputs); err != nil {
			return 0, nil, err
		} else if entity == nil {
			return 0, nil, HttpError(404, n.ID, "not found.")
		} else {
			return 200, entity, nil
		}
	}
}

func (n *Node) innerGET(parent *ResolvedNode, id string) (interface{}, error) {
	if n.IsCollection {
		return n.innerGETCollection(parent, id)
	}
	entity, _, err := n.Ops["Manifest"].Invoke(bindManifestInputs(n.Parent, id))
	if err != nil {
		return nil, err
	} else if entity == nil {
		return nil, Error(id, "not found")
	} else {
		return entity, nil
	}
}

func (n *Node) innerGETCollection(parent *ResolvedNode, id string) (interface{}, error) {
	entity, other, err := n.Ops["Page"].Invoke(bindManifestInputs(n.Parent, id))
	if err != nil {
		return nil, err
	} else if ids, ok := other.([]string); !ok {
		return nil, Error("page must return a slice of strings and and error")
	} else if entity == nil {
		return nil, Error("page must set the receiver to a non-nil collection")
	} else {
		collection := map[string]interface{}{}
		rn := &ResolvedNode{n, parent, id, entity, nil}
		for _, id := range ids {
			if cn, err := rn.Locate(id); err != nil {
				collection[id] = err
			} else if cr, err := cn.MemberResource(n.CollectionTag.EmbedFields); err != nil {
				collection[id] = err
			} else {
				collection[id] = cr
			}
		}
		return collection, nil
	}
}

func bindManifestInputs(n *Node, id string) map[IN]boundInput {
	var nilParent interface{}
	if n != nil {
		nilParent = reflect.New(n.EntityType).Interface()
	} else {
		nilParent = nil
	}
	return map[IN]boundInput{
		IN_Parent: func(_ *BoundOp) (interface{}, error) {
			return nilParent, nil
		},
		IN_ID: func(_ *BoundOp) (interface{}, error) {
			return id, nil
		},
		IN_PageNum: func(_ *BoundOp) (interface{}, error) {
			return 1, nil
		},
	}
}
