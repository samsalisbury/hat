package hat

import (
	"reflect"
)

func (n *Node) AnalyseOperations() error {
	ops := map[string]*Operation{}
	for name, op := range OPS {
		if m, ok := n.EntityPtrType.MethodByName(name); !ok {
			continue
		} else err := op.AssertMethodConforms(m); err != nil {
			return nil, Error()
		}
	}
	if _, ok := ops["Manifest"]; !ok {
		return Error(n.EntityType, "does not have a Manifest method.")
	}
	return nil
}

func (n *Node) AnalyseMembers() error {
	return nil
}

func (n *Node) AnalyseCollection() error () {
	return nil
}

func (n *Node) AnalyseHttpMethods() error () {
	return nil
}