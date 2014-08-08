package hat

type Member struct {
	Node             *Node
	DefaultExpansion string
}

func newMember(n *Node, defaultExpansion string) *Member {
	return &Member{n, defaultExpansion}
}
