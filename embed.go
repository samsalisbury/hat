package hat

import (
	"reflect"
	"strings"
)

type smap map[string]interface{}

func (n *ResolvedNode) EmbedResources() (embedded smap, entity smap, err error) {
	entity, err = mapify(n.Entity)
	if err != nil {
		return nil, nil, err
	}
	if n.Node.IsCollection {
		return embedItems(&entity, n)
	}
	return embedMembers(&entity, n)
}

func embedItems(entity *smap, n *ResolvedNode) (smap, smap, error) {
	items := stringMap(n.Entity)
	embedded := make(map[string]interface{}, len(items))
	for id, _ := range items {
		if childNode, err := n.Resolve(id); err != nil {
			return nil, nil, err
		} else if embeddedResource, err := childNode.Resource(nil); err != nil {
			return nil, nil, err
		} else {
			embedded[id] = embeddedResource
			deleteIgnoringCase(entity, id)
		}
	}
	return embedded, *entity, nil
}

func embedMembers(entity *smap, n *ResolvedNode) (smap, smap, error) {
	embedded := smap{}
	for name, member := range n.Node.Members {
		if !member.Tag.Embed {
			continue
		}
		if memberNode, err := n.Locate(name); err != nil {
			return nil, nil, err
		} else if resource, err := memberNode.MemberResource(member.Tag.EmbedFields); err != nil {
			return nil, nil, err
		} else {
			embedded[name] = resource
			deleteIgnoringCase(entity, name)
		}
	}
	return embedded, *entity, nil
}

func deleteIgnoringCase(from *smap, k string) {
	for realKey, _ := range *from {
		if strings.EqualFold(k, realKey) {
			delete(*from, realKey)
		}
	}
}

func stringMap(collection interface{}) smap {
	v := reflect.ValueOf(collection).Elem()
	m := smap{}
	for _, kv := range v.MapKeys() {
		m[kv.String()] = v.MapIndex(kv).Interface()
	}
	return m
}
