package hat

func (n *ResolvedNode) Embeds() (smap, error) {
	if n.Node.IsCollection {
		return embedItems(n)
	}
	return embedMembers(n)
}

func embedItems(n *ResolvedNode) (smap, error) {
	embedded := make(smap, len(n.CollectionIDs))
	for _, id := range n.CollectionIDs {
		if childNode, err := n.Resolve(id); err != nil {
			return nil, err
		} else if resource, err := childNode.MemberResource(n.Node.Collection.Tag.EmbedFields); err != nil {
			return nil, err
		} else {
			embedded[id] = resource
		}
	}
	return embedded, nil
}

func embedMembers(n *ResolvedNode) (smap, error) {
	embedded := smap{}
	for urlName, member := range n.Node.Members {
		if !member.Tag.Embed {
			continue
		}
		if memberNode, err := n.Locate(urlName); err != nil {
			return nil, err
		} else if resource, err := memberNode.MemberResource(member.Tag.EmbedFields); err != nil {
			return nil, err
		} else {
			embedded[member.Name] = resource
		}
	}
	return embedded, nil
}
