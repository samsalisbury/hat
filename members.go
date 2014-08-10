package hat

type Member struct {
	Node *Node
	Tag  *Tag
}

func newMember(n *Node, tag *Tag) (*Member, error) {
	return &Member{n, tag}, nil
}
