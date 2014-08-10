package hat

func (n *Node) initOps() error {
	n.Ops = map[string]*CompiledOp{}
	for name, op := range op_specs {
		if m, ok := n.EntityPtrType.MethodByName(name); !ok {
			continue
		} else if co, err := op.Compile(n, m); err != nil {
			return err
		} else {
			n.Ops[name] = co
		}
	}
	if n.Ops["Manifest"] == nil {
		return Error(n.EntityType, "does not have a Manifest method.")
	}
	return nil
}
