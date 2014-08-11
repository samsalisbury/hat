package hat

import (
	"encoding/json"
	"reflect"
	"strings"
)

type smap map[string]interface{}

// This is only necessary because the toSmap method below
// uses JSON encoder to convert to smap, which means
func (m *smap) deleteIgnoringCase(k string) {
	for realKey, _ := range *m {
		if strings.EqualFold(k, realKey) {
			delete(*m, realKey)
		}
	}
}

// Converts the collection to a smap or panics.
func collectionToSmap(collection interface{}) smap {
	v := reflect.ValueOf(collection).Elem()
	m := smap{}
	for _, kv := range v.MapKeys() {
		m[kv.String()] = v.MapIndex(kv).Interface()
	}
	return m
}

// THIS IS VERY BAD. FIX IT!
// Converts anything to smap. Probably not a very good
// idea, since JSON encoding has too many other
// rules.
// TODO: Consider an alternative method (Gob? Protobuf? Hand-rolled?)
func toSmap(v interface{}) (smap, error) {
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
