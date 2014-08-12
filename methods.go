package hat

import "reflect"

type StdHTTPMethod func() (statusCode int, resource *Resource, err error)

func makeHTTPMethods(n ResolvedNode, inputs map[IN]boundInput) map[string]StdHTTPMethod {
	return map[string]StdHTTPMethod{
		"GET": makeGET(n, inputs),
	}
}

func makeGET(n ResolvedNode, inputs map[IN]boundInput) StdHTTPMethod {
	return func() (statusCode int, resource *Resource, err error) {
		if r, err := n.Resource(); err != nil {
			debug("Resource error:", err)
			return 0, nil, err
		} else {
			debug("Resource success...", r)
			return 200, r, nil
		}
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
