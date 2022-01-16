package activity

import "database-camp/internal/utils"

type Choices interface {
	CreatePropositionChoices() interface{}
}

type MultipleChoice struct {
	ID        int    `gorm:"primaryKey;column:multiple_choice_id" json:"multiple_choice_id"`
	Content   string `gorm:"column:content" json:"content"`
	IsCorrect bool   `gorm:"column:is_correct" json:"is_correct"`
}

type MultipleChoices []MultipleChoice

func (choices MultipleChoices) CreatePropositionChoices() interface{} {
	preparedChoices := make([]map[string]interface{}, 0)

	utils.Shuffle(choices)

	for _, v := range choices {
		preparedChoice, _ := utils.StructToMap(v)
		delete(preparedChoice, "is_correct")
		preparedChoices = append(preparedChoices, preparedChoice)
	}

	return preparedChoices
}

type CompletionChoice struct {
	ID            int    `gorm:"primaryKey;column:completion_choice_id" json:"completion_choice_id"`
	Content       string `gorm:"column:content" json:"content"`
	QuestionFirst string `gorm:"column:question_first" json:"question_first"`
	QuestionLast  string `gorm:"column:question_last" json:"question_last"`
}

type CompletionChoices []CompletionChoice

func (choices CompletionChoices) CreatePropositionChoices() interface{} {
	contents := make([]interface{}, 0)
	questions := make([]interface{}, 0)

	for _, v := range choices {
		contents = append(contents, v.Content)
		questions = append(questions, map[string]interface{}{
			"id":    v.ID,
			"first": v.QuestionFirst,
			"last":  v.QuestionLast,
		})
	}

	utils.Shuffle(contents)
	utils.Shuffle(questions)

	prepared := map[string]interface{}{
		"contents":  contents,
		"questions": questions,
	}

	return prepared
}

type MatchingChoice struct {
	ID        int    `gorm:"primaryKey;column:matching_choice_id" json:"matching_choice_id"`
	PairItem1 string `gorm:"column:pair_item1" json:"pair_item1"`
	PairItem2 string `gorm:"column:pair_item2" json:"pair_item2"`
}

type MatchingChoices []MatchingChoice

func (choices MatchingChoices) CreatePropositionChoices() interface{} {
	pairItem1List := make([]interface{}, 0)
	pairItem2List := make([]interface{}, 0)

	for _, v := range choices {
		pairItem1List = append(pairItem1List, v.PairItem1)
		pairItem2List = append(pairItem2List, v.PairItem2)
	}

	utils.Shuffle(pairItem1List)
	utils.Shuffle(pairItem2List)

	prepared := map[string]interface{}{
		"items_left":  pairItem1List,
		"items_right": pairItem2List,
	}

	return prepared
}
