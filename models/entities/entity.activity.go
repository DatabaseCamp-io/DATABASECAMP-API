package entities

// entity.activity.go
/**
 * 	This file is a part of models, used to collect model for entities of activity
 */

import (
	"DatabaseCamp/errs"
	"DatabaseCamp/models/request"
	"DatabaseCamp/models/storages"
	"DatabaseCamp/utils"
)

// Model of hint roadmap result for activity
type HintRoadMap struct {
	Level       int `json:"level"`
	ReducePoint int `json:"reduce_point"`
}

// Model of activity hint result for Exam activity
type ActivityHint struct {
	TotalHint   int               `json:"total_hint"`
	UsedHints   []storages.HintDB `json:"used_hints"`
	HintRoadMap []HintRoadMap     `json:"hint_roadmap"`
}

// Model of activity detail result for activity
type ActivityDetail struct {
	ID        int    `json:"activity_id"`
	TypeID    int    `json:"activity_type_id"`
	ContentID *int   `json:"content_id"`
	Order     int    `json:"activity_order"`
	Story     string `json:"story"`
	Point     int    `json:"point"`
	Question  string `json:"question"`
}

/**
 * This class manage activity model
 */
type Activity struct {
	info               ActivityDetail
	propositionChoices interface{}
	choices            interface{}
	hint               *ActivityHint
}

/**
 * Getter for getting hint
 *
 * @return hint
 */
func (a *Activity) GetHint() *ActivityHint {
	return a.hint
}

/**
 * Getter for getting proposition choices of the activity
 *
 * @return propositionChoices of the activity
 */
func (a *Activity) GetPropositionChoices() interface{} {
	return a.propositionChoices
}

/**
 * Getter for getting information of the activity
 *
 * @return information of the activity
 */
func (a *Activity) GetInfo() ActivityDetail {
	return a.info
}

/**
 * Setter for set activity
 *
 * @param activityDB to set activity data
 */
func (a *Activity) SetActivity(activityDB storages.ActivityDB) {
	utils.NewType().StructToStruct(activityDB, &a.info)
}

/**
 * Setter for set choice by choiceDB
 *
 * @param choiceDB to set choice
 *
 * @return the error of setting
 */
func (a *Activity) SetChoicesByChoiceDB(choiceDB interface{}) error {
	a.choices = choiceDB
	if a.info.TypeID == 1 {
		propositionChoices, choiceOK := choiceDB.([]storages.MatchingChoiceDB)
		if !choiceOK {
			return errs.NewInternalServerError("รูปแบบของตัวเลือกไม่ถูกต้อง", "Invalid Choice")
		}
		a.propositionChoices = a.PrepareMatchingChoice(propositionChoices)
	} else if a.info.TypeID == 2 {
		propositionChoices, choiceOK := choiceDB.([]storages.MultipleChoiceDB)
		if !choiceOK {
			return errs.NewInternalServerError("รูปแบบของตัวเลือกไม่ถูกต้อง", "Invalid Choice")
		}
		a.propositionChoices = a.PrepareMultipleChoice(propositionChoices)
	} else if a.info.TypeID == 3 {
		propositionChoices, choiceOK := choiceDB.([]storages.CompletionChoiceDB)
		if !choiceOK {
			return errs.NewInternalServerError("รูปแบบของตัวเลือกไม่ถูกต้อง", "Invalid Choice")
		}
		a.propositionChoices = a.PrepareCompletionChoice(propositionChoices)
	} else {
		a.propositionChoices = nil
	}
	return nil
}

/**
 * Setter for set hint
 *
 * @param activityHints Activity hints model from database to set hint
 * @param userHintsDB 	User hints model from database to set hint
 */
func (a *Activity) SetHint(activityHints []storages.HintDB, userHintsDB []storages.UserHintDB) {
	a.hint = &ActivityHint{
		TotalHint: len(activityHints),
	}
	for _, hint := range activityHints {
		if a.isUsedHint(userHintsDB, hint.ID) {
			a.hint.UsedHints = append(a.hint.UsedHints, hint)
		}
		a.hint.HintRoadMap = append(a.hint.HintRoadMap, HintRoadMap{
			Level:       hint.Level,
			ReducePoint: hint.PointReduce,
		})
	}
}

/**
 * Check if hint is used
 *
 * @param userHintsDB	Hints of the user to check
 * @param hintID		Hint ID to check
 *
 * @return true if hint is used, false if not
 */
func (a *Activity) isUsedHint(userHintsDB []storages.UserHintDB, hintID int) bool {
	for _, userHint := range userHintsDB {
		if userHint.HintID == hintID {
			return true
		}
	}
	return false
}

/**
 * Prepare multipleChoice
 *
 * @param 	multipleChoice 	Multiple choice to prepare
 *
 * @return prepared multiple choice
 */
func (a *Activity) PrepareMultipleChoice(multipleChoice []storages.MultipleChoiceDB) interface{} {
	preparedChoices := make([]map[string]interface{}, 0)
	utils.NewHelper().Shuffle(multipleChoice)
	for _, v := range multipleChoice {
		preparedChoice, _ := utils.NewType().StructToMap(v)
		delete(preparedChoice, "is_correct")
		preparedChoices = append(preparedChoices, preparedChoice)
	}

	return preparedChoices
}

/**
 * Prepare matchingchoice
 *
 * @param 	matchingChoice 	Matching choice to prepare
 *
 * @return prepared matching choice
 */
func (a *Activity) PrepareMatchingChoice(matchingChoice []storages.MatchingChoiceDB) interface{} {
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

/**
 * Prepare completionchoice
 *
 * @param 	completionchoice 	Completion choice to prepare
 *
 * @return prepared completion choice
 */
func (a *Activity) PrepareCompletionChoice(completionChoice []storages.CompletionChoiceDB) interface{} {
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

/**
 * Check answer for matching choice
 *
 * @param 	choices 	Choices of activity
 * @param 	answer 		Answer of matching choice
 *
 * @return 	true if input answer is correct, false if not
 */
func (a *Activity) IsMatchingCorrect(choices []storages.MatchingChoiceDB, answer []request.PairItemRequest) bool {
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

/**
 * Check answer for completion choice
 *
 * @param 	choices 	Choices of activity
 * @param 	answer 		Answer of completion choice
 *
 * @return 	true if input answer is correct, false if not
 */
func (a *Activity) IsCompletionCorrect(choices []storages.CompletionChoiceDB, answer []request.PairContentRequest) bool {
	for _, correct := range choices {
		for _, answer := range answer {
			if (correct.ID == *answer.ID) && (correct.Content != *answer.Content) {
				return false
			}
		}
	}
	return true
}

/**
 * Check answer for multiple choice
 *
 * @param 	choices 	Choices of activity
 * @param 	answer 		Answer of multiple choice
 *
 * @return 	true if input answer is correct, false if not
 */
func (a *Activity) IsMultipleCorrect(choices []storages.MultipleChoiceDB, answer int) bool {
	for _, v := range choices {
		if v.IsCorrect && v.ID == answer {
			return true
		}
	}
	return false
}

/**
 * Convert type of choice to pair item
 *
 * @param 	raw 	Item for convert to pair item
 *
 * @return  choice that converted to pair item
 * @return  the error converting
 */
func (a *Activity) convertToPairItem(raw interface{}) ([]request.PairItemRequest, error) {
	result := make([]request.PairItemRequest, 0)
	list, ok := raw.([]interface{})
	if !ok {
		return nil, errs.NewBadRequestError("รูปแบบของคำตอบไม่ถูกต้อง", "Invalid Answer Format")
	}
	for _, v := range list {
		temp := request.PairItemRequest{}
		utils.NewType().StructToStruct(v, &temp)
		err := temp.Validate()
		if err != nil {
			return nil, err
		}
		result = append(result, temp)
	}
	return result, nil
}

/**
 * Check answer for matching choice
 *
 * @param 	answer 	Answer of activity
 *
 * @return 	true if input answer is correct, false if not
 * @return  the error of checking matching choice
 */
func (a *Activity) checkMatchingCorrect(answer interface{}) (bool, error) {
	matchingChoices, choiceOK := a.choices.([]storages.MatchingChoiceDB)
	_answer, answerOK := answer.([]request.PairItemRequest)
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

/**
 * Check answer for multiple choice
 *
 * @param 	answer 	answer of activity
 *
 * @return 	true if input answer is correct, false if not
 * @return  the error of checking multiple choice
 */
func (a *Activity) checkMultipleCorrect(answer interface{}) (bool, error) {
	multipleChoices, choiceOK := a.choices.([]storages.MultipleChoiceDB)
	if !choiceOK {
		return false, errs.NewBadRequestError("รูปแบบของคำตอบไม่ถูกต้อง", "Invalid Answer Format")
	}
	return a.IsMultipleCorrect(multipleChoices, utils.NewType().ParseInt(answer)), nil
}

/**
 * Convert type of content to pair content
 *
 * @param 	raw 		item for convert to pair content request
 *
 * @return  choice that converted to pair content
 * @return  the error of converting
 */
func (a *Activity) convertToPairContent(raw interface{}) ([]request.PairContentRequest, error) {
	result := make([]request.PairContentRequest, 0)
	list, ok := raw.([]interface{})
	if !ok {
		return nil, errs.NewBadRequestError("รูปแบบของคำตอบไม่ถูกต้อง", "Invalid Answer Format")
	}
	for _, v := range list {
		temp := request.PairContentRequest{}
		utils.NewType().StructToStruct(v, &temp)
		err := temp.Validate()
		if err != nil {
			return nil, err
		}
		result = append(result, temp)
	}
	return result, nil
}

/**
 * Check answer for completion choice
 *
 * @param 	answer 	Answer of activity
 *
 * @return 	true if input answer is correct, false if not
 * @return  the error of checking completion choice
 */
func (a *Activity) checkCompletionCorrect(answer interface{}) (bool, error) {
	completionChoices, choiceOK := a.choices.([]storages.CompletionChoiceDB)
	_answer, answerOK := answer.([]request.PairContentRequest)
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

/**
 * Check if answer is correct
 *
 * @param 	answer 	Answer of activity
 *
 * @return 	true if the answer is correct, false if not
 * @return  the error of checking answer
 */
func (a *Activity) IsAnswerCorrect(answer interface{}) (bool, error) {
	if a.info.TypeID == 1 {
		return a.checkMatchingCorrect(answer)
	} else if a.info.TypeID == 2 {
		return a.checkMultipleCorrect(answer)
	} else if a.info.TypeID == 3 {
		return a.checkCompletionCorrect(answer)
	} else {
		return false, errs.NewBadRequestError("ประเภทของกิจกรรมไม่ถูกต้อง", "Invalid Activity Type")
	}
}
