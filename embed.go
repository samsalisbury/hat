package hat

func (n *ResolvedNode) EmbeddedResources() (map[string]interface{}, error) {
	embedded := map[string]interface{}{}
	for name, member := range n.Node.Members {
		if !member.Tag.Embed {
			continue
		}
		if memberNode, err := n.Locate(name); err != nil {
			return nil, err
		} else if resolvedMemberNode, err := memberNode.Node.Resolve(n, n.ID); err != nil {
			return nil, err
		} else {
			if resource, err := resolvedMemberNode.MemberResource(member.Tag.EmbedFields); err != nil {
				return nil, err
			} else {
				embedded[name] = resource
			}
		}
	}
	return embedded, nil
}
