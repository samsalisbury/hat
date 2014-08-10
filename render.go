package hat

import (
	"strings"
)

type inBinder func(*ResolvedNode) map[IN]boundInput

func (root *Node) Render(path string, method string, inputBinder inBinder) (int, interface{}, error) {
	target, err := root.Locate(strings.Split(path[1:], "/")...)
	if err != nil {
		return 0, nil, err
	}
	target.HTTPMethods = makeHTTPMethods(target)
	inputs := inputBinder(target)
	if method, ok := target.HTTPMethods[method]; !ok {
		return 0, nil, HttpError(405, target, "does not support method", method, "; it does support:", supportedMethods(target))
	} else if statusCode, entity, err := method(inputs); err != nil {
		return 0, nil, err
	} else if resource, err := target.Resource(entity); err != nil {
		return 0, nil, err
	} else {
		return statusCode, resource, nil
	}
}

func supportedMethods(n *ResolvedNode) string {
	m := []string{}
	for k, _ := range n.HTTPMethods {
		m = append(m, k)
	}
	return strings.Join(m, ", ")
}
