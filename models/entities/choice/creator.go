package choice

import (
	"DatabaseCamp/models/storages"
	"errors"
)

type Choice interface {
	CreatePropositionChoices() PropositionChoices
}

type creator struct {
}

type multipleChoice struct {
	multipleChoiceDB   storages.MultipleChoiceDB
	propositionChoices multiplePropositionChoices
}

func NewCreator() *creator {
	return &creator{}
}

func (c multipleChoice) CreatePropositionChoices() PropositionChoices {
	return c.propositionChoices
}

func (creator) CreateChoice(choiceDB interface{}) (Choice, error) {
	multipleChoiceDB, ok := choiceDB.(storages.MultipleChoiceDB)
	if !ok {
		return multipleChoice{}, errors.New("")
	}
	return multipleChoice{multipleChoiceDB: multipleChoiceDB}, nil
}
