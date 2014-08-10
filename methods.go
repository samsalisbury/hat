package hat

type StdHTTPMethod func(inputs map[IN]boundInput) (statusCode int, resource interface{}, err error)

func makeHTTPMethods(n *ResolvedNode) map[string]StdHTTPMethod {
	return map[string]StdHTTPMethod{
		"GET": makeGET(n),
	}
}

func makeGET(n *ResolvedNode) StdHTTPMethod {
	return func(inputs map[IN]boundInput) (statusCode int, resource interface{}, err error) {
		if resource, _, err := n.Node.Ops["Manifest"].Invoke(inputs); err != nil {
			return 0, nil, err
		} else if resource == nil {
			return 404, nil, Error(n.ID, "not found.")
		} else {
			return 200, resource, nil
		}
	}
}

func (n *Node) innerGET(inputs map[IN]boundInput) (resource interface{}, err error) {
	if resource, _, err := n.Ops["Manifest"].Invoke(inputs); err != nil {
		return nil, err
	} else if resource == nil {
		return nil, Error("not found")
	} else {
		return resource, nil
	}
}
