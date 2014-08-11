package hat

import (
	"strconv"
	"strings"
)

type Tag struct {
	Embed       bool
	EmbedFields []string
	Link        bool
	LinkRel     string
	Page        bool
	PageNum     int
	PageSize    int
}

func parseTag(datum string) (*Tag, error) {
	data := strings.Split(datum, ";")
	tag := &Tag{}
	for _, t := range data {
		if err := parseTagDirective(t, tag); err != nil {
			return nil, err
		}
	}
	return tag, nil
}

func parseTagDirective(data string, tag *Tag) error {
	data = strings.Trim(data, " \t")
	data = strings.TrimSuffix(data, ")")
	data = strings.Trim(data, " \t")
	parts := strings.SplitN(data, "(", 2)
	if len(parts) != 2 {
		return Error("Tag", data, "not recognised. Format is tagname(params)")
	}

	if fn, ok := tagMap[parts[0]]; !ok {
		return Error("Tag name", parts[0], "not recognised expected link or embed")
	} else if err := fn(parts[1], tag); err != nil {
		return err
	}
	return nil
}

var tagMap = map[string]func(string, *Tag) error{
	"embed": embedTag,
	"link":  linkTag,
	"page":  pageTag,
}

func embedTag(params string, tag *Tag) error {
	fields := []string{}
	if len(params) != 0 {
		fields = strings.Split(params, ",")
	}
	tag.Embed = true
	tag.EmbedFields = fields
	return nil
}

func linkTag(params string, tag *Tag) error {
	tag.Link = true
	tag.LinkRel = params
	return nil
}

func pageTag(params string, tag *Tag) error {
	tag.Page = true
	parts := strings.Split(params, ",")
	if len(parts) > 2 {
		return Error("page tag; got", params, "expected 2 digits separated with comma")
	}
	pageNum := parts[0]
	var page, size int
	if len(pageNum) == 0 {
		page = 0
	} else if p, err := strconv.ParseInt(pageNum, 10, 32); err != nil {
		return Error("Unable to parse page number:", quot(pageNum))
	} else {
		page = int(p)
	}
	tag.PageNum = page
	if len(parts) == 1 {
		return nil
	}
	pageSize := parts[1]
	if len(pageSize) == 0 {
		return nil
	} else if s, err := strconv.ParseInt(pageSize, 10, 32); err != nil {
		return Error("Unable to parse page size:", quot(pageSize))
	} else {
		size = int(s)
	}
	tag.PageSize = size
	return nil
}
