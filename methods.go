package hat

type StdHTTPMethod func(n *Node, parent interface{}, id string, p *Payload) (statusCode int, resource interface{}, err error)

func makeGET(n *Node) StdHTTPMethod {
	return func(n *Node, parent interface{}, id string, p *Payload) (statusCode int, resource interface{}, err error) {
		println(":::: GET CALLED WITH Parent ::: " + Error(parent).Error())
		if resource, _, err := n.Operations.Manifest.Invoke(n, parent, id, p); err != nil {
			return 0, nil, err
		} else if resource == nil {
			return 404, nil, Error(id, "not found.")
		} else {
			return 200, resource, nil
		}
	}
}
