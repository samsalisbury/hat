package hat

import (
	"reflect"
)

type Operation struct {
	On      SELF
	Inputs  []IN
	Outputs []OUT
}

type SELF int

const (
	SELF_Nil        = IN(iota)
	SELF_Payload    = IN(iota)
	SELF_Manifested = IN(iota)
)

type IN int

const (
	IN_Parent      = IN(iota)
	IN_OtherEntity = IN(iota)
	IN_ID          = IN(iota)
)

type OUT int

const (
	OUT_Error       = OUT(iota)
	OUT_OtherEntity = OUT(iota)
)

func on(self SELF) *Operation                      { return &Operation{Self: self} }
func (o *Operation) In(inputs ...IN) *Operation    { o.Inputs = inputs; return o }
func (o *Operation) Out(outputs ...OUT) *Operation { o.Outputs = outputs; return o }

var OPS = map[string]*Operation{
	"Manifest": on(SELF_Nil).In(IN_Entity_Nil, IN_Parent, IN_ID).Out(OUT_Error),
}

func (op *Operation) AssertMethodConforms(t reflect.Type) error {
	return nil
}
