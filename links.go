package hat

func (n *ResolvedNode) Links(entity *smap) ([]Link, error) {
	links := []Link{Link{"self", n.Path()}}
	for name, member := range n.Node.Members {
		if member.Tag.Link {
			rel := member.Tag.LinkRel
			if rel == "" {
				rel = name
			}
			links = append(links, Link{rel, n.Path() + "/" + name})
			entity.deleteIgnoringCase(name)
		}
	}
	return links, nil
}

type Link struct {
	Rel   string `json:"rel"`
	Hrefs string `json:"href"`
}
