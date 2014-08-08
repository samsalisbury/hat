package hat

import (
	"encoding/json"
)

type Resource struct {
	Entity   interface{}
	Embedded map[string]interface{}
	Links    map[string]interface{}
}

func (n *ResolvedNode) Resource(entity interface{}) (*Resource, error) {
	if embedded, err := n.ManifestEmbedded(n.Node.Members); err != nil {
		return nil, err
	} else if links, err := n.ManifestLinks(); err != nil {
		return nil, err
	} else {
		return &Resource{entity, embedded, links}, nil
	}
}

func (n *ResolvedNode) ManifestEmbedded(members map[string]*Member) (map[string]interface{}, error) {
	return nil, nil
}

func (n *ResolvedNode) ManifestLinks() (map[string]interface{}, error) {
	return nil, nil
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
