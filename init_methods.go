package hat

func (n *Node) initHttpMethods() error {
	n.HTTPMethods = map[string]StdHTTPMethod{}
	if n.Operations.Manifest != nil {
		n.HTTPMethods["GET"] = makeGET(n)
	}
	return nil
}
