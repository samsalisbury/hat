package hat

import (
	"strings"
)

func (root *Node) Render(path string, method string, p *Payload) (int, interface{}, error) {
	if target, err := root.Locate(strings.Split(path[1:], "/")...); err != nil {
		return 0, nil, err
	} else if method, ok := target.Node.HTTPMethods[method]; !ok {
		return 0, nil, HttpError(405, target, "does not support method", method, "; it does support:", supportedMethods(target.Node))
	} else if statusCode, entity, err := method(target.Node, target.ParentEntity(), target.ID, p); err != nil {
		return 0, nil, err
	} else if resource, err := target.Resource(entity); err != nil {
		return 0, nil, err
	} else {
		return statusCode, resource, nil
	}
}

func supportedMethods(n *Node) string {
	m := []string{}
	for k, _ := range n.HTTPMethods {
		m = append(m, k)
	}
	return strings.Join(m, ", ")
}
