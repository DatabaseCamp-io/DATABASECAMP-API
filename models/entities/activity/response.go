package activity

import (
	"DatabaseCamp/models/entities/choice"
	"DatabaseCamp/models/storages"
)

type hintRoadmap struct {
	Level       int `json:"level"`
	ReducePoint int `json:"reduce_point"`
}

type hint struct {
	TotalHint   int               `json:"total_hint"`
	UsedHints   []storages.HintDB `json:"used_hints"`
	HintRoadMap []hintRoadmap     `json:"hint_roadmap"`
}

type Response struct {
	Activity storages.ActivityDB       `json:"activity"`
	Choices  choice.PropositionChoices `json:"choice"`
	Hint     hint                      `json:"hint"`
}
