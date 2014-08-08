package hat

import (
	"reflect"
)

func (n *Node) initOperations() error {
	ops := &compiled_ops{}
	numDefinedOps := op_specs_V.NumField()
	for i := 0; i < numDefinedOps; i++ {
		field := op_specs_V.Field(i)
		name := op_specs_T.Field(i).Name
		op := field.Interface().(*Operation)
		if m, ok := n.EntityPtrType.MethodByName(name); !ok {
			continue
		} else if compiled, err := op.Compile(n, m); err != nil {
			return err
		} else {
			f := reflect.ValueOf(ops).Elem().FieldByName(name)
			f.Set(reflect.ValueOf(compiled))
		}
	}
	if ops.Manifest == nil {
		return Error(n.EntityType, "does not have a Manifest method.")
	}
	n.Operations = *ops
	return nil
}
