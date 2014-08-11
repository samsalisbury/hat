package hat

import "reflect"

type StdHTTPMethod func() (statusCode int, resource interface{}, err error)

func makeHTTPMethods(n *ResolvedNode, inputs map[IN]boundInput) map[string]StdHTTPMethod {
	return map[string]StdHTTPMethod{
		"GET": makeGET(n, inputs),
	}
}

func makeGET(n *ResolvedNode, inputs map[IN]boundInput) StdHTTPMethod {
	return func() (statusCode int, other interface{}, err error) {
		var sc int
		if n.Entity == nil {
			sc = 404
		} else {
			sc = 200
		}
		return sc, nil, nil
	}
	// if n.Node.IsCollection {
	// 	return makeGETCollection(n, inputs)
	// } else {
	// 	return makeGETSingular(n, inputs)
	// }
}

func makeGETCollection(n *ResolvedNode, inputs map[IN]boundInput) StdHTTPMethod {
	return func() (statusCode int, other interface{}, err error) {
		// GET is a special case, since all resources are "GOT" before other things happen to them.
		var sc int
		if n.Entity == nil {
			sc = 404
		} else {
			sc = 200
		}
		return sc, nil, nil
	}
}

func makeGETSingular(n *ResolvedNode, inputs map[IN]boundInput) StdHTTPMethod {
	return func() (statusCode int, other interface{}, err error) {
		return 200, nil, nil
	}
}

func (parentNode *Node) createChildManifestInputs(parentEntity interface{}, id string) map[IN]boundInput {
	if parentEntity == nil && parentNode != nil {
		parentEntity = reflect.New(parentNode.EntityType).Interface()
	}
	return map[IN]boundInput{
		IN_Parent: func(_ *BoundOp) (interface{}, error) {
			return parentEntity, nil
		},
		IN_ID: func(_ *BoundOp) (interface{}, error) {
			return id, nil
		},
		IN_PageNum: func(_ *BoundOp) (interface{}, error) {
			return 1, nil
		},
	}
}
