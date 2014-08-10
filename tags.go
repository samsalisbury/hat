package hat

import (
	"strings"
)

type Tag struct {
	Embed       bool
	EmbedFields []string
	Link        bool
	LinkRel     string
}

func parseTag(tag string) (*Tag, error) {
	tag = strings.Trim(tag, " \t")
	tag = strings.TrimSuffix(tag, ")")
	tag = strings.Trim(tag, " \t")
	parts := strings.SplitN(tag, "(", 2)
	if len(parts) != 2 {
		return nil, Error("Tag", tag, "not recognised. Format is tagname(params)")
	}
	switch parts[0] {
	default:
		return nil, Error("Tag name", tag, "not recognised expected link or embed")
	case "embed":
		return embedTag(parts[1])
	case "link":
		return linkTag(parts[1])
	}
}

func embedTag(params string) (*Tag, error) {
	fields := []string{}
	if len(params) != 0 {
		fields = strings.Split(params, ",")
	}
	return &Tag{Embed: true, EmbedFields: fields}, nil
}

func linkTag(params string) (*Tag, error) {
	return &Tag{Link: true, LinkRel: params}, nil
}
