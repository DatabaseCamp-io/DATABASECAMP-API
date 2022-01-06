package choice

import (
	"reflect"
)

type PropositionChoices interface {
	getPropositionChoicesType() reflect.Type
}

type multiplePropositionChoices []map[string]interface{}

func (c multiplePropositionChoices) getPropositionChoicesType() reflect.Type {
	return reflect.TypeOf(c)
}
