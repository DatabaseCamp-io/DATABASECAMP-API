package models

import (
	"DatabaseCamp/errs"
	"DatabaseCamp/utils"
)

type hintRoadMap struct {
	Level       int `json:"level"`
	ReducePoint int `json:"reduce_point"`
}

type activityHint struct {
	TotalHint   int           `json:"total_hint"`
	UsedHints   []HintDB      `json:"used_hints"`
	HintRoadMap []hintRoadMap `json:"hint_roadmap"`
}

type activityDetail struct {
	ID        int    `json:"activity_id"`
	TypeID    int    `json:"activity_type_id"`
	ContentID *int   `json:"content_id"`
	Order     int    `json:"activity_order"`
	Story     string `json:"story"`
	Point     int    `json:"point"`
	Question  string `json:"question"`
}

type activity struct {
	Info               activityDetail `json:"activity"`
	PropositionChoices interface{}    `json:"proposition_choices"`
	Choices            interface{}    `json:"choices"`
	Hint               *activityHint  `json:"hint"`
}

func NewActivity() *activity {
	return &activity{}
}

func (a *activity) ToPropositionResponse() *ActivityResponse {
	response := ActivityResponse{
		Activity: a.Info,
		Choices:  a.PropositionChoices,
		Hint:     *a.Hint,
	}
	return &response
}

func (a *activity) ToAnswerResponse(updatedPoint int, isCorrect bool) *AnswerResponse {
	response := AnswerResponse{
		ActivityID:   a.Info.ID,
		IsCorrect:    isCorrect,
		UpdatedPoint: updatedPoint,
	}
	return &response
}

func (a *activity) PrepareActivity(activityDB ActivityDB) {
	utils.NewType().StructToStruct(activityDB, &a.Info)
}

func (a *activity) PrepareChoicesByChoiceDB(choiceDB interface{}) error {
	a.Choices = choiceDB
	if a.Info.TypeID == 1 {
		propositionChoices, choiceOK := choiceDB.([]MatchingChoiceDB)
		if !choiceOK {
			return errs.NewInternalServerError("รูปแบบของตัวเลือกไม่ถูกต้อง", "Invalid Choice")
		}
		a.PropositionChoices = a.PrepareMatchingChoice(propositionChoices)
	} else if a.Info.TypeID == 2 {
		propositionChoices, choiceOK := choiceDB.([]MultipleChoiceDB)
		if !choiceOK {
			return errs.NewInternalServerError("รูปแบบของตัวเลือกไม่ถูกต้อง", "Invalid Choice")
		}
		a.PropositionChoices = a.PrepareMultipleChoice(propositionChoices)
	} else if a.Info.TypeID == 3 {
		propositionChoices, choiceOK := choiceDB.([]CompletionChoiceDB)
		if !choiceOK {
			return errs.NewInternalServerError("รูปแบบของตัวเลือกไม่ถูกต้อง", "Invalid Choice")
		}
		a.PropositionChoices = a.PrepareCompletionChoice(propositionChoices)
	} else {
		a.PropositionChoices = nil
	}
	return nil
}

func (a *activity) PrepareHint(activityHints []HintDB, userHintsDB []UserHintDB) {
	a.Hint = &activityHint{
		TotalHint: len(activityHints),
	}
	for _, hint := range activityHints {
		if a.isUsedHint(userHintsDB, hint.ID) {
			a.Hint.UsedHints = append(a.Hint.UsedHints, hint)
		}
		a.Hint.HintRoadMap = append(a.Hint.HintRoadMap, hintRoadMap{
			Level:       hint.Level,
			ReducePoint: hint.PointReduce,
		})
	}
}

func (a *activity) isUsedHint(userHintsDB []UserHintDB, hintID int) bool {
	for _, userHint := range userHintsDB {
		if userHint.HintID == hintID {
			return true
		}
	}
	return false
}

func (a *activity) PrepareMultipleChoice(multipleChoice []MultipleChoiceDB) interface{} {
	preparedChoices := make([]map[string]interface{}, 0)
	utils.NewHelper().Shuffle(multipleChoice)
	for _, v := range multipleChoice {
		preparedChoice, _ := utils.NewType().StructToMap(v)
		delete(preparedChoice, "is_correct")
		preparedChoices = append(preparedChoices, preparedChoice)
	}

	return preparedChoices
}

func (a *activity) PrepareMatchingChoice(matchingChoice []MatchingChoiceDB) interface{} {
	pairItem1List := make([]interface{}, 0)
	pairItem2List := make([]interface{}, 0)
	for _, v := range matchingChoice {
		pairItem1List = append(pairItem1List, v.PairItem1)
		pairItem2List = append(pairItem2List, v.PairItem2)
	}
	utils.NewHelper().Shuffle(pairItem1List)
	utils.NewHelper().Shuffle(pairItem2List)
	prepared := map[string]interface{}{
		"items_left":  pairItem1List,
		"items_right": pairItem2List,
	}
	return prepared
}

func (a *activity) PrepareCompletionChoice(completionChoice []CompletionChoiceDB) interface{} {
	contents := make([]interface{}, 0)
	questions := make([]interface{}, 0)
	for _, v := range completionChoice {
		contents = append(contents, v.Content)
		questions = append(questions, map[string]interface{}{
			"id":    v.ID,
			"first": v.QuestionFirst,
			"last":  v.QuestionLast,
		})
	}
	utils.NewHelper().Shuffle(contents)
	utils.NewHelper().Shuffle(questions)
	prepared := map[string]interface{}{
		"contents":  contents,
		"questions": questions,
	}
	return prepared
}

func (a *activity) IsMatchingCorrect(choices []MatchingChoiceDB, answer []PairItemRequest) bool {
	Item1Item2Map := map[string]string{}
	for _, correct := range choices {
		Item1Item2Map[correct.PairItem1] = correct.PairItem2
	}
	for _, answer := range answer {
		if Item1Item2Map[*answer.Item1] != *answer.Item2 && Item1Item2Map[*answer.Item2] != *answer.Item1 {
			return false
		}
	}
	return true
}

func (a *activity) IsCompletionCorrect(choices []CompletionChoiceDB, answer []PairContentRequest) bool {
	for _, correct := range choices {
		for _, answer := range answer {
			if (correct.ID == *answer.ID) && (correct.Content != *answer.Content) {
				return false
			}
		}
	}
	return true
}

func (a *activity) IsMultipleCorrect(choices []MultipleChoiceDB, answer int) bool {
	for _, v := range choices {
		if v.IsCorrect && v.ID == answer {
			return true
		}
	}
	return false
}

func (a *activity) convertToPairItem(raw interface{}) ([]PairItemRequest, error) {
	result := make([]PairItemRequest, 0)
	list, ok := raw.([]interface{})
	if !ok {
		return nil, errs.NewBadRequestError("รูปแบบของคำตอบไม่ถูกต้อง", "Invalid Answer Format")
	}
	for _, v := range list {
		temp := PairItemRequest{}
		utils.NewType().StructToStruct(v, &temp)
		err := temp.Validate()
		if err != nil {
			return nil, err
		}
		result = append(result, temp)
	}
	return result, nil
}

func (a *activity) checkMatchingCorrect(answer interface{}) (bool, error) {
	matchingChoices, choiceOK := a.Choices.([]MatchingChoiceDB)
	_answer, answerOK := answer.([]PairItemRequest)
	if !answerOK {
		var err error
		_answer, err = a.convertToPairItem(answer)
		if err != nil {
			return false, err
		}
	}
	if len(matchingChoices) != len(_answer) || !choiceOK {
		return false, errs.NewBadRequestError("รูปแบบของคำตอบไม่ถูกต้อง", "Invalid Answer Format")
	}
	return a.IsMatchingCorrect(matchingChoices, _answer), nil
}

func (a *activity) checkMultipleCorrect(answer interface{}) (bool, error) {
	multipleChoices, choiceOK := a.Choices.([]MultipleChoiceDB)
	if !choiceOK {
		return false, errs.NewBadRequestError("รูปแบบของคำตอบไม่ถูกต้อง", "Invalid Answer Format")
	}
	return a.IsMultipleCorrect(multipleChoices, utils.NewType().ParseInt(answer)), nil
}

func (a *activity) convertToPairContent(raw interface{}) ([]PairContentRequest, error) {
	result := make([]PairContentRequest, 0)
	list, ok := raw.([]interface{})
	if !ok {
		return nil, errs.NewBadRequestError("รูปแบบของคำตอบไม่ถูกต้อง", "Invalid Answer Format")
	}
	for _, v := range list {
		temp := PairContentRequest{}
		utils.NewType().StructToStruct(v, &temp)
		err := temp.Validate()
		if err != nil {
			return nil, err
		}
		result = append(result, temp)
	}
	return result, nil
}

func (a *activity) checkCompletionCorrect(answer interface{}) (bool, error) {
	completionChoices, choiceOK := a.Choices.([]CompletionChoiceDB)
	_answer, answerOK := answer.([]PairContentRequest)
	if !answerOK || !choiceOK {
		var err error
		_answer, err = a.convertToPairContent(answer)
		if err != nil {
			return false, err
		}
	}
	if len(completionChoices) != len(_answer) || !choiceOK {
		return false, errs.NewBadRequestError("รูปแบบของคำตอบไม่ถูกต้อง", "Invalid Answer Format")
	}
	return a.IsCompletionCorrect(completionChoices, _answer), nil
}

func (a *activity) IsAnswerCorrect(answer interface{}) (bool, error) {
	if a.Info.TypeID == 1 {
		return a.checkMatchingCorrect(answer)
	} else if a.Info.TypeID == 2 {
		return a.checkMultipleCorrect(answer)
	} else if a.Info.TypeID == 3 {
		return a.checkCompletionCorrect(answer)
	} else {
		return false, errs.NewBadRequestError("ประเภทของกิจกรรมไม่ถูกต้อง", "Invalid Activity Type")
	}
}
