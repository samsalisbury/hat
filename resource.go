package hat

import (
	"encoding/json"
)

type Resource struct {
	Entity   interface{}
	Embedded map[string]interface{}
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
	m, err := mapify(n.Entity)
	if err != nil {
		return nil, err
	}
	filtered := make(map[string]interface{}, len(fields))
	for _, f := range fields {
		filtered[f] = m[f]
	}
	return n.Resource(filtered)
}

func (n *ResolvedNode) Resource(other interface{}) (*Resource, error) {
	if other != nil {
		// Other could be anything, e.g. a status message,
		// therefore it can't be linked like other entities.
		// However I think we need a way to allow adding arbitrary
		// links to these kinds of responses.
		return &Resource{other, nil, nil}, nil
	}
	if embedded, entity, err := n.EmbedResources(); err != nil {
		return nil, err
	} else if links, err := n.Links(); err != nil {
		return nil, err
	} else {
		return &Resource{entity, embedded, links}, nil
	}
}

func mapify(v interface{}) (map[string]interface{}, error) {
	if j, err := json.Marshal(v); err != nil {
		return nil, err
	} else {
		var m map[string]interface{}
		if err := json.Unmarshal(j, &m); err != nil {
			return nil, err
		}
		return m, nil
	}
}
