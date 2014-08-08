package hat

func (n *ResolvedNode) Links() ([]Link, error) {
	return []Link{Link{"self", n.Path()}}, nil
}

type Link struct {
	Rel   string `json:"rel"`
	Hrefs string `json:"href"`
}
