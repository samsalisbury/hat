package hat

type Resource struct {
	Entity   interface{}
	Embedded smap
	Links    []Link
}

func (n *ResolvedNode) MemberResource(fields []string) (*Resource, error) {
	if len(fields) != 0 {
		return n.FilteredMemberResource(fields)
	}
	return n.DefaultMemberResource()
}

func (n *ResolvedNode) DefaultMemberResource() (*Resource, error) {
	return n.Resource(n.Entity)
}

func (n *ResolvedNode) FilteredMemberResource(fields []string) (*Resource, error) {
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

func (n *ResolvedNode) Resource(other interface{}) (*Resource, error) {
	if other != nil {
		// Other could be anything, e.g. a status message,
		// therefore it can't be linked like other entities.
		// However I think we need a way to allow adding arbitrary
		// links to these kinds of responses.
		return &Resource{other, nil, nil}, nil
	}
	if entity, err := toSmap(n.Entity); err != nil {
		return nil, err
	} else if embedded, err := n.Embeds(); err != nil {
		return nil, err
	} else if links, err := n.Links(); err != nil {
		return nil, err
	} else {
		for k, _ := range embedded {
			entity.deleteIgnoringCase(k)
		}
		for _, l := range links {
			entity.deleteIgnoringCase(l.Rel)
		}
		return &Resource{entity, embedded, links}, nil
	}
}
