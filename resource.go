package hat

type Resource struct {
	Entity                  interface{}
	EmbeddedMembers         map[string]*Resource
	EmbeddedCollectionItems []*Resource
	Links                   []Link
}
