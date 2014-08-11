package hat

func (n *ResolvedNode) EmbedResources() (embedded, entity smap, err error) {
	entity, err = toSmap(n.Entity)
	if err != nil {
		return nil, nil, err
	}
	if n.Node.IsCollection {
		embedded, err = embedItems(&entity, n)
	} else {
		embedded, err = embedMembers(&entity, n)
	}
	return embedded, entity, err
}

func embedItems(entity *smap, n *ResolvedNode) (smap, error) {
	items := collectionToSmap(n.Entity)
	embedded := make(smap, len(items))
	for id, _ := range items {
		if childNode, err := n.Resolve(id); err != nil {
			return nil, err
		} else if embeddedResource, err := childNode.Resource(nil); err != nil {
			return nil, err
		} else {
			embedded[id] = embeddedResource
			entity.deleteIgnoringCase(id)
		}
	}
	return embedded, nil
}

func embedMembers(entity *smap, n *ResolvedNode) (smap, error) {
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
			entity.deleteIgnoringCase(member.Name)
		}
	}
	return embedded, nil
}
