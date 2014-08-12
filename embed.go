package hat

func (n *ResolvedNode) EmbeddedCollectionItems() ([]*Resource, error) {
	if !n.Node.IsCollection {
		return nil, nil
	}
	items := make([]*Resource, len(n.CollectionIDs))
	for i, id := range n.CollectionIDs {
		if childNode, err := n.Resolve(id); err != nil {
			return nil, err
		} else if resource, err := childNode.EmbeddedResource(n.Node.Collection.Tag); err != nil {
			return nil, err
		} else {
			items[i] = resource
		}
	}
	return items, nil
}

func (n *ResolvedNode) EmbeddedMembers() (map[string]*Resource, error) {
	if n.Node.IsCollection {
		return nil, nil
	}
	embedded := map[string]*Resource{}
	for urlName, member := range n.Node.Members {
		if !member.Tag.Embed {
			continue
		}
		if memberNode, err := n.Locate(urlName); err != nil {
			return nil, err
		} else if resource, err := memberNode.EmbeddedResource(member.Tag); err != nil {
			return nil, err
		} else {
			embedded[member.Name] = resource
		}
	}
	return embedded, nil
}
