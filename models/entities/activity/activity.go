package activity

import (
	"DatabaseCamp/models/entities/choice"
	"DatabaseCamp/models/storages"
)

type userHints []storages.UserHintDB

func (h userHints) IsUsed(hintID int) bool {
	for _, userHint := range h {
		if userHint.HintID == hintID {
			return true
		}
	}
	return false
}

type Activity interface {
	NewResponse() *Response
}

type activity struct {
	activityDB      *storages.ActivityDB
	userHintsDB     []storages.UserHintDB
	activityHintsDB []storages.HintDB
	choiceDB        interface{}
}

func New(
	activityDB *storages.ActivityDB,
	userHintsDB []storages.UserHintDB,
	activityHintsDB []storages.HintDB,
	choiceDB interface{},
) *activity {
	return &activity{
		activityDB:      activityDB,
		userHintsDB:     userHintsDB,
		activityHintsDB: activityHintsDB,
		choiceDB:        choiceDB,
	}
}

func (a activity) NewResponse() (*Response, error) {
	choices, err := a.prepareChoices()

	return &Response{
		Activity: *a.activityDB,
		Choices:  choices,
		Hint:     a.prepareHintInfo(),
	}, err
}

func (a activity) prepareChoices() (choice.PropositionChoices, error) {
	choiceCreator := choice.NewCreator()

	choice, err := choiceCreator.CreateChoice(a.choiceDB)
	if err != nil {
		return nil, err
	}

	return choice.CreatePropositionChoices(), nil
}

func (a activity) prepareHintInfo() hint {
	userHints := userHints(a.userHintsDB)
	usedHints := make([]storages.HintDB, 0)
	roadmap := make([]hintRoadmap, 0)

	for _, hint := range a.activityHintsDB {

		if userHints.IsUsed(hint.ID) {
			usedHints = append(usedHints, hint)
		}

		roadmap = append(roadmap, hintRoadmap{
			Level:       hint.Level,
			ReducePoint: hint.PointReduce,
		})
	}

	return hint{
		TotalHint:   len(a.activityHintsDB),
		UsedHints:   usedHints,
		HintRoadMap: roadmap,
	}
}
