package hat

import (
	"reflect"
)

var op_specs = ops{
	Manifest: on(SELF_Nil).In(IN_Parent, IN_ID).Out(OUT_Error),
}

var op_specs_V = reflect.ValueOf(op_specs)

type ops struct {
	Manifest *Operation
}

type compiled_ops struct {
	Manifest *CompiledOperation
}

var op_specs_T = reflect.TypeOf(compiled_ops{})

type Operation struct {
	On      SELF
	Inputs  []IN
	Outputs []OUT
}

func (o *Operation) RequiresNilReceiver() bool {
	return o.On == SELF_Nil
}

func (o *Operation) RequiresPayloadReceiver() bool {
	return o.On == SELF_Payload
}

func (o *Operation) RequiresManifestedReceiver() bool {
	return o.On == SELF_Manifested
}

func (o *Operation) RequiresOtherPayload() bool {
	return o.Requires(IN_OtherPayload)
}

func (o *Operation) Requires(in IN) bool {
	for _, i := range o.Inputs {
		if i == in {
			return true
		}
	}
	return false
}

type CompiledOperation struct {
	Def             *Operation
	OtherEntityType reflect.Type
	Method          reflect.Method
}

func (o *Operation) Compile(n *Node, m reflect.Method) (*CompiledOperation, error) {
	co := &CompiledOperation{o, nil, m}
	actualNumIn := m.Type.NumIn() - 1
	if actualNumIn != len(o.Inputs) {
		return nil, n.MethodError(m.Name, "Wrong number of inputs. Expected", len(o.Inputs), "but got", actualNumIn)
	}
	for i, in := range o.Inputs {
		realType := m.Type.In(i + 1)
		if err := in.Accepts(n, m.Name, i, realType); err != nil {
			return nil, err
		}
		switch in {
		case IN_OtherPayload:
			co.OtherEntityType = realType
		default:
		}
	}
	return co, nil
}

type BoundOperation struct {
	Compiled     *CompiledOperation
	Receiver     interface{}
	OtherPayload interface{}
	Method       reflect.Value
}

func (co *CompiledOperation) BindNilReceiver(t reflect.Type) *BoundOperation {
	return &BoundOperation{Compiled: co, Receiver: reflect.New(t).Interface()}
}

func (co *CompiledOperation) BindManifestedOrPayloadReciever(rcvr interface{}) (*BoundOperation, error) {
	return &BoundOperation{Compiled: co, Receiver: rcvr}, nil
}

func (co *CompiledOperation) Invoke(n *Node, parent interface{}, id string, p *Payload) (entity interface{}, other interface{}, err error) {
	if bo, err := co.Bind(n, parent, id, p); err != nil {
		return nil, nil, err
	} else {
		return bo.Invoke(n, parent, id)
	}
}

func (bo *BoundOperation) Invoke(n *Node, parent interface{}, id string) (entity interface{}, other interface{}, err error) {
	in := bo.PrepareInputs(n, parent, id)
	out := bo.Method.Call(in)
	entity = bo.Receiver
	for i, o := range bo.Compiled.Def.Outputs {
		if o == OUT_Error {
			if !out[i].IsNil() {
				err = out[i].Interface().(error)
			}
		} else if o == OUT_OtherEntity {
			if !out[i].IsNil() {
				other = out[i].Interface()
			}
		}
	}
	return entity, other, err
}

func (bo *BoundOperation) PrepareInputs(n *Node, parent interface{}, id string) []reflect.Value {
	prepared := []reflect.Value{}
	if bo.Compiled.Def.Requires(IN_Parent) {
		// TODO: Evaluate if this check is necessary
		if parent == nil {
			if n.Parent != nil {
				parent = reflect.New(n.Parent.EntityType).Interface()
			} else {
				parent = struct{}{}
			}
		}
		prepared = append(prepared, reflect.ValueOf(parent))
	}
	if bo.Compiled.Def.Requires(IN_ID) {
		prepared = append(prepared, reflect.ValueOf(id))
	}
	if bo.Compiled.Def.Requires(IN_OtherPayload) {
		prepared = append(prepared, reflect.ValueOf(bo.OtherPayload))
	}
	return prepared
}

func (co *CompiledOperation) Bind(n *Node, parent interface{}, id string, p *Payload) (*BoundOperation, error) {
	if bo, err := co.BindReceiver(n, parent, id, p); err != nil {
		return nil, err
	} else if err := bo.BindMethod(); err != nil {
		return nil, err
	} else if err := bo.BindOtherPayload(n, p); err != nil {
		return nil, err
	} else {
		return bo, nil
	}
}

func (bo *BoundOperation) BindOtherPayload(n *Node, p *Payload) error {
	if !bo.Compiled.Def.Requires(IN_OtherPayload) {
		return nil
	}
	if payload, err := p.Manifest(bo.Compiled.OtherEntityType); err != nil {
		return err
	} else {
		bo.OtherPayload = payload
	}
	return nil
}

func (bo *BoundOperation) BindMethod() error {
	method := reflect.ValueOf(bo.Receiver).MethodByName(bo.Compiled.Method.Name)
	bo.Method = method
	return nil
}

func (co *CompiledOperation) BindReceiver(n *Node, parent interface{}, id string, p *Payload) (*BoundOperation, error) {
	if co.Def.RequiresNilReceiver() {
		return co.BindNilReceiver(n.EntityType), nil
	} else if co.Def.RequiresManifestedReceiver() {
		if rcvr, _, err := n.Operations.Manifest.Invoke(n, parent, id, p); err != nil {
			return nil, err
		} else {
			return co.BindManifestedOrPayloadReciever(rcvr)
		}
	} else if co.Def.RequiresPayloadReceiver() {
		if rcvr, err := p.Manifest(n.EntityType); err != nil {
			return nil, err
		} else {
			return co.BindManifestedOrPayloadReciever(rcvr)
		}
	} else {
		panic("No receiver type specified for *" + n.EntityType.Name() + "." + co.Method.Name)
	}
}

type SELF int

const (
	SELF_Nil        = SELF(iota)
	SELF_Payload    = SELF(iota)
	SELF_Manifested = SELF(iota)
)

type IN int

const (
	IN_Parent       = IN(iota)
	IN_OtherPayload = IN(iota)
	IN_ID           = IN(iota)
)

func (in IN) Accepts(n *Node, name string, pos int, t reflect.Type) error {
	switch in {
	default:
		panic("The programmer has made a serious error.")
	case IN_ID:
		if t.Kind() == reflect.String {
			return nil
		} else {
			return n.MethodError(name, "cannot accept input type", t, "at position", pos)
		}
	case IN_Parent:
		if n.Parent == nil || t == n.Parent.EntityPtrType {
			return nil // maybe one day we won't need these useless params
		} else {
			return n.MethodError(name, "expects a pointer to its parent type", n.Parent.EntityPtrType, "at position", pos)
		}
	case IN_OtherPayload:
		if t.Kind() == reflect.Ptr {
			elemKind := t.Elem().Kind()
			switch elemKind {
			case reflect.Struct, reflect.Map, reflect.Slice:
				return nil // This is the only ok case, otherwise we return the below error.
			}
		}
		return n.MethodError(name, "expects a pointer to a struct, map, or slice at position", pos)
	}
}

type OUT int

const (
	OUT_Error       = OUT(iota)
	OUT_OtherEntity = OUT(iota)
)

func on(self SELF) *Operation                      { return &Operation{On: self} }
func (o *Operation) In(inputs ...IN) *Operation    { o.Inputs = inputs; return o }
func (o *Operation) Out(outputs ...OUT) *Operation { o.Outputs = outputs; return o }

func iType(nilPtr interface{}) reflect.Type {
	return reflect.TypeOf(nilPtr).Elem()
}

func typ(example interface{}) reflect.Type {
	return reflect.TypeOf(example)
}
