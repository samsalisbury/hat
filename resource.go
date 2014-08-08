package hat

import (
	"encoding/json"
)

type Resource struct {
	Entity   interface{}
	Embedded map[string]interface{}
	Links    []Link
}

func (n *ResolvedNode) DefaultResource() (*Resource, error) {
	return n.Resource(n.Entity)
}

func (n *ResolvedNode) Resource(entity interface{}) (*Resource, error) {
	if embedded, err := n.EmbeddedResources(); err != nil {
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
