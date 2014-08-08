package hat

import (
	"strings"
)

type Member struct {
	Node *Node
	Tag  *MemberTag
}

type MemberTag struct {
	Embed   bool
	Link    bool
	LinkRel string
}

func parseMemberTag(tag string) (*MemberTag, error) {
	tag = strings.Trim(tag, " \t")
	tag = strings.TrimSuffix(tag, ")")
	tag = strings.Trim(tag, " \t")
	parts := strings.SplitN(tag, "(", 2)
	if len(parts) != 2 {
		return nil, Error("Tag", tag, "not recognised. Should be either link(...) or embed(...)")
	} else {
		if parts[0] == "embed" {
			return memberEmbedTag(parts[1])
		} else if parts[0] == "link" {
			return memberLinkTag(parts[1])
		} else {
			return nil, Error("Tag", tag, "not recognised. Should be either link(...) or embed(...)")
		}
	}
}

func memberEmbedTag(params string) (*MemberTag, error) {
	return &MemberTag{Embed: true}, nil
}

func memberLinkTag(params string) (*MemberTag, error) {
	return &MemberTag{Link: true, LinkRel: params}, nil
}

func newMember(n *Node, tag string) (*Member, error) {
	if memberTag, err := parseMemberTag(tag); err != nil {
		return nil, err
	} else {
		return &Member{n, memberTag}, nil
	}
}
