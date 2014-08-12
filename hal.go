package hat

func halError(err error) hatError {
	return Error("RenderAsHAL:", err)
}

func RenderAsHAL(r *Resource) (smap, error) {
	hal := smap{}
	e, err := toSmap(r.Entity)
	if err != nil {
		return nil, halError(err)
	}
	for k, v := range e {
		hal[k] = v
	}
	if r.Links != nil && len(r.Links) != 0 {
		if links, err := halLinks(r.Links); err != nil {
			return nil, halError(err)
		} else {
			hal["_links"] = links
		}
	}
	if r.EmbeddedMembers != nil && len(r.EmbeddedMembers) != 0 {
		embedded := smap{}
		for name, childResource := range r.EmbeddedMembers {
			if c, err := RenderAsHAL(childResource); err != nil {
				return nil, halError(err)
			} else {
				embedded[name] = c
			}
		}
		hal["_embedded"] = embedded
	} else if r.EmbeddedCollectionItems != nil && len(r.EmbeddedCollectionItems) != 0 {
		embedded := make([]smap, len(r.EmbeddedCollectionItems))
		for i, itemResource := range r.EmbeddedCollectionItems {
			if c, err := RenderAsHAL(itemResource); err != nil {
				return nil, halError(err)
			} else {
				embedded[i] = c
			}
		}
		hal["_embedded"] = smap{"NAME": embedded}
	}
	return hal, nil
}

func halLinks(links []Link) (smap, error) {
	raw := rawHalLinks(links)
	flattened := make(smap, len(raw))
	for rel, hrefs := range raw {
		if len(hrefs) == 1 {
			flattened[rel] = map[string]string{"href": hrefs[0]}
		} else {
			hh := make([]map[string]string, len(hrefs))
			for i, h := range hrefs {
				hh[i] = map[string]string{"href": h}
			}
			flattened[rel] = hh
		}
	}
	return flattened, nil
}

func rawHalLinks(links []Link) map[string][]string {
	raw := map[string][]string{}
	for _, l := range links {
		if _, ok := raw[l.Rel]; ok {
			raw[l.Rel] = append(raw[l.Rel], l.Href)
		} else {
			raw[l.Rel] = []string{l.Href}
		}
	}
	return raw
}

func halEmbeddedMembers(members map[string]*Resource) (smap, error) {
	return nil, nil
}

func halEmbeddedCollectionItems(members []*Resource) (map[string][]interface{}, error) {
	return nil, nil
}
