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
		pos := i + 1
		realType := m.Type.In(pos)
		if realType.Kind() == reflect.Ptr {
			realType = realType.Elem()
		}
		if !in.Accepts(realType) {
			return nil, n.Error(m.Name, "cannot accept input type", realType, "at position", i)
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
	println(Error("Invoke bo.Method =", bo.Method).Error())
	println(Error("Invoke bo.Method (NumIn) =", bo.Method.Type().NumIn()).Error())
	in := bo.PrepareInputs(n, parent, id)
	for i, j := range in {
		println(Error("In :", i, "=", j).Error())
	}
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
	println(Error("PrepareInputs parent:", parent, "; id:", id).Error())
	prepared := []reflect.Value{}
	if bo.Compiled.Def.Requires(IN_Parent) {
		println("ADD PARENT")
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
		println("ADD ID")
		prepared = append(prepared, reflect.ValueOf(id))
	}
	if bo.Compiled.Def.Requires(IN_OtherPayload) {
		println("ADD OTHER PAYLOAD")
		prepared = append(prepared, reflect.ValueOf(bo.OtherPayload))
	}
	println(Error("PREPRED", len(prepared), "ARGS").Error())
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
	println(Error("BindMethod bo.Receiver =", bo.Receiver).Error())
	println(Error("bo.Compiled.Method.Name =", bo.Compiled.Method.Name).Error())
	println(Error("reflect.ValueOf(bo.Receiver) =", reflect.ValueOf(bo.Receiver)).Error())
	println(Error("bo.Receiver =", bo.Receiver).Error())
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

func (in IN) Accepts(t reflect.Type) bool {
	if in == IN_ID {
		println("Checking IN_ID", Error(t).Error())
		return t.Kind() == reflect.String
	} else {
		println("Checking Other IN", Error(t).Error())
		return true
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
