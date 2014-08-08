package hat

func (n *ResolvedNode) EmbeddedResources() (map[string]interface{}, error) {
	embedded := map[string]interface{}{}
	for name, member := range n.Node.Members {
		expansion := member.DefaultExpansion
		if expansion == "full()" {
			memberNode, _ := n.Locate(name)
			resolvedMemberNode, _ := memberNode.Node.Resolve(n, n.ID)
			resource, _ := resolvedMemberNode.DefaultResource()
			embedded[name] = resource
		}
	}
	return embedded, nil
}
