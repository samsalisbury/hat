package hat

type Resource struct {
	Entity                  interface{}
	EmbeddedMembers         map[string]*Resource
	EmbeddedCollectionItems []*Resource
	Links                   []Link
}

func (n *ResolvedNode) Resource(other interface{}) (*Resource, error) {
	if other != nil {
		// Other could be anything, e.g. a status message,
		// therefore it can't be linked like other entities.
		// However I think we need a way to allow adding arbitrary
		// links to these kinds of responses.
		return &Resource{other, nil, nil, nil}, nil
	}
	if entity, err := toSmap(n.Entity); err != nil {
		return nil, err
	} else if embeddedMembers, err := n.EmbeddedMembers(); err != nil {
		return nil, err
	} else if embeddedCollectionItems, err := n.EmbeddedCollectionItems(); err != nil {
		return nil, err
	} else if links, err := n.Links(); err != nil {
		return nil, err
	} else {
		for k, _ := range embeddedMembers {
			entity.deleteIgnoringCase(k)
		}
		for _, l := range links {
			entity.deleteIgnoringCase(l.Rel)
		}
		return &Resource{entity, embeddedMembers, embeddedCollectionItems, links}, nil
	}
}

func (n *ResolvedNode) EmbeddedResource(tag *Tag) (*Resource, error) {
	if len(tag.EmbedFields) != 0 {
		return n.FilteredEmbeddedResource(tag.EmbedFields)
	}
	return n.DefaultEmbeddedResource()
}

func (n *ResolvedNode) DefaultEmbeddedResource() (*Resource, error) {
	return n.Resource(n.Entity)
}

func (n *ResolvedNode) FilteredEmbeddedResource(fields []string) (*Resource, error) {
	m, err := toSmap(n.Entity)
	if err != nil {
		return nil, err
	}
	filtered := make(smap, len(fields))
	for _, f := range fields {
		filtered[f] = m[f]
	}
	return n.Resource(nil)
}
