package hat

import "reflect"

type StdHTTPMethod func() (statusCode int, resource interface{}, err error)

func makeHTTPMethods(n *ResolvedNode, inputs map[IN]boundInput) map[string]StdHTTPMethod {
	return map[string]StdHTTPMethod{
		"GET": makeGET(n, inputs),
	}
}

func makeGET(n *ResolvedNode, inputs map[IN]boundInput) StdHTTPMethod {
	if n.Node.IsCollection {
		return makeGETCollection(n, inputs)
	} else {
		return makeGETMember(n, inputs)
	}
}

func makeGETCollection(n *ResolvedNode, inputs map[IN]boundInput) StdHTTPMethod {
	return func() (statusCode int, entity interface{}, err error) {
		return 200, n.Entity, nil
	}
}

func makeGETMember(n *ResolvedNode, inputs map[IN]boundInput) StdHTTPMethod {
	return func() (statusCode int, entity interface{}, err error) {
		return 200, n.Entity, nil
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
