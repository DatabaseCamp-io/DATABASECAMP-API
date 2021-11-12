package services

import (
	"DatabaseCamp/errs"
	"DatabaseCamp/models"
	"DatabaseCamp/utils"
)

type activityManager struct {
}

type IActivityManager interface {
	PrepareMultipleChoice(multipleChoice []models.MultipleChoiceDB) interface{}
	PrepareMatchingChoice(matchingChoice []models.MatchingChoiceDB) interface{}
	PrepareCompletionChoice(completionChoice []models.CompletionChoiceDB) interface{}
	IsMatchingCorrect(choices []models.MatchingChoiceDB, answer []models.PairItemRequest)
	IsCompletionCorrect(choices []models.CompletionChoiceDB, answer []models.PairContentRequest)
	IsMultipleCorrect(choices []models.MultipleChoiceDB, answer int)
	IsAnswerCorrect(typeID int, choice interface{}, answer interface{}) (bool, error)
}

func NewActivityManager() *activityManager {
	return &activityManager{}
}

func (m activityManager) PrepareMultipleChoice(multipleChoice []models.MultipleChoiceDB) interface{} {
	preparedChoices := make([]map[string]interface{}, 0)
	utils.NewHelper().Shuffle(multipleChoice)
	for _, v := range multipleChoice {
		preparedChoice, _ := utils.NewType().StructToMap(v)
		delete(preparedChoice, "is_correct")
		preparedChoices = append(preparedChoices, preparedChoice)
	}

	return preparedChoices
}

func (m activityManager) PrepareMatchingChoice(matchingChoice []models.MatchingChoiceDB) interface{} {
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

func (m activityManager) PrepareCompletionChoice(completionChoice []models.CompletionChoiceDB) interface{} {
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

func (m activityManager) IsMatchingCorrect(choices []models.MatchingChoiceDB, answer []models.PairItemRequest) bool {
	for _, correct := range choices {
		for _, answer := range answer {
			if (correct.PairItem1 == *answer.Item1) && (correct.PairItem2 != *answer.Item2) {
				return false
			}
		}
	}
	return true
}

func (m activityManager) IsCompletionCorrect(choices []models.CompletionChoiceDB, answer []models.PairContentRequest) bool {
	for _, correct := range choices {
		for _, answer := range answer {
			if (correct.ID == *answer.ID) && (correct.Content != *answer.Content) {
				return false
			}
		}
	}
	return true
}

func (m activityManager) IsMultipleCorrect(choices []models.MultipleChoiceDB, answer int) bool {
	for _, v := range choices {
		if v.IsCorrect && v.ID == answer {
			return true
		}
	}
	return false
}

func (m activityManager) convertToPairItem(raw interface{}) ([]models.PairItemRequest, error) {
	result := make([]models.PairItemRequest, 0)
	list, ok := raw.([]interface{})
	if !ok {
		return nil, errs.NewBadRequestError("รูปแบบของคำตอบไม่ถูกต้อง", "Invalid Answer Format")
	}
	for _, v := range list {
		temp := models.PairItemRequest{}
		utils.NewType().StructToStruct(v, &temp)
		err := temp.Validate()
		if err != nil {
			return nil, err
		}
		result = append(result, temp)
	}
	return result, nil
}

func (m activityManager) checkMatchingCorrect(choice interface{}, answer interface{}) (bool, error) {
	matchingChoices, choiceOK := choice.([]models.MatchingChoiceDB)
	_answer, answerOK := answer.([]models.PairItemRequest)
	if !answerOK {
		var err error
		_answer, err = m.convertToPairItem(answer)
		if err != nil {
			return false, err
		}
	}
	if len(matchingChoices) != len(_answer) || !choiceOK {
		return false, errs.NewBadRequestError("รูปแบบของคำตอบไม่ถูกต้อง", "Invalid Answer Format")
	}
	return m.IsMatchingCorrect(matchingChoices, _answer), nil
}

func (m activityManager) checkMultipleCorrect(choice interface{}, answer interface{}) (bool, error) {
	multipleChoices, choiceOK := choice.([]models.MultipleChoiceDB)
	if !choiceOK {
		return false, errs.NewBadRequestError("รูปแบบของคำตอบไม่ถูกต้อง", "Invalid Answer Format")
	}
	return m.IsMultipleCorrect(multipleChoices, utils.NewType().ParseInt(answer)), nil
}

func (m activityManager) convertToPairContent(raw interface{}) ([]models.PairContentRequest, error) {
	result := make([]models.PairContentRequest, 0)
	list, ok := raw.([]interface{})
	if !ok {
		return nil, errs.NewBadRequestError("รูปแบบของคำตอบไม่ถูกต้อง", "Invalid Answer Format")
	}
	for _, v := range list {
		temp := models.PairContentRequest{}
		utils.NewType().StructToStruct(v, &temp)
		err := temp.Validate()
		if err != nil {
			return nil, err
		}
		result = append(result, temp)
	}
	return result, nil
}

func (m activityManager) checkCompletionCorrect(choice interface{}, answer interface{}) (bool, error) {
	completionChoices, choiceOK := choice.([]models.CompletionChoiceDB)
	_answer, answerOK := answer.([]models.PairContentRequest)
	if !answerOK || !choiceOK {
		var err error
		_answer, err = m.convertToPairContent(answer)
		if err != nil {
			return false, err
		}
	}
	if len(completionChoices) != len(_answer) {
		return false, errs.NewBadRequestError("รูปแบบของคำตอบไม่ถูกต้อง", "Invalid Answer Format")
	}
	return m.IsCompletionCorrect(completionChoices, _answer), nil
}

func (m activityManager) IsAnswerCorrect(typeID int, choice interface{}, answer interface{}) (bool, error) {
	if typeID == 1 {
		return m.checkMatchingCorrect(choice, answer)
	} else if typeID == 2 {
		return m.checkMultipleCorrect(choice, answer)
	} else if typeID == 3 {
		return m.checkCompletionCorrect(choice, answer)
	} else {
		return false, errs.NewBadRequestError("ประเภทของกิจกรรมไม่ถูกต้อง", "Invalid Activity Type")
	}
}
