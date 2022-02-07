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
	countCorrect := 0
	preparedChoices := make([]map[string]interface{}, 0)

	utils.Shuffle(choices)

	for _, v := range choices {
		if v.IsCorrect {
			countCorrect++
		}

		preparedChoice, _ := utils.StructToMap(v)
		delete(preparedChoice, "is_correct")
		preparedChoices = append(preparedChoices, preparedChoice)
	}

	isMultipleAnswers := countCorrect > 1

	result := map[string]interface{}{
		"is_multiple_answers": isMultipleAnswers,
		"choices":             preparedChoices,
	}

	return result
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

type VocalGroupChoice struct {
	GroupName string `gorm:"column:vocab_group_name" json:"vocab_group_name"`
	Vocab     string `gorm:"column:vocab" json:"vocab"`
}

type VocalGroupChoices []VocalGroupChoice

func (choices VocalGroupChoices) CreatePropositionChoices() interface{} {
	groups := make([]string, 0)
	vocabs := make([]string, 0)

	utils.Shuffle(choices)

	for _, v := range choices {
		groups = append(groups, v.GroupName)
		vocabs = append(vocabs, v.Vocab)
	}

	preparedChoices := map[string]interface{}{
		"groups": groups,
		"vocabs": vocabs,
	}

	return preparedChoices
}

type Determinant struct {
	Value string `gorm:"column:value" json:"value"`
	Fixed bool   `gorm:"column:fixed" json:"-"`
}

type Dependency struct {
	Dependent    string        `gorm:"column:dependent" json:"dependent"`
	Fixed        bool          `gorm:"column:fixed" json:"-"`
	Determinants []Determinant `gorm:"foreignKey:determinant_id" json:"determinants"`
}

type DependencyChoice struct {
	ID           int          `gorm:"column:dependency_choice_id"`
	Dependencies []Dependency `gorm:"foreignKey:dependency_id"`
}

func (choice DependencyChoice) CreatePropositionChoices() interface{} {

	type dependency struct {
		Dependent         *string  `json:"dependent"`
		DeterminantsCount int      `json:"determinants_count"`
		Determinants      []string `json:"determinants"`
	}

	type result struct {
		Vocabs       []string     `json:"vocabs"`
		Dependencies []dependency `json:"dependencies"`
	}

	propositionChoices := result{
		Vocabs:       make([]string, 0),
		Dependencies: make([]dependency, 0),
	}

	for _, v := range choice.Dependencies {
		dependencyResult := dependency{
			DeterminantsCount: len(v.Determinants),
			Determinants:      make([]string, 0),
		}

		if v.Fixed {
			dependencyResult.Dependent = &v.Dependent
		} else {
			propositionChoices.Vocabs = append(propositionChoices.Vocabs, v.Dependent)
		}

		for _, d := range v.Determinants {
			if d.Fixed {
				dependencyResult.Determinants = append(dependencyResult.Determinants, d.Value)
			} else {
				propositionChoices.Vocabs = append(propositionChoices.Vocabs, d.Value)
			}
		}

		propositionChoices.Dependencies = append(propositionChoices.Dependencies, dependencyResult)
	}

	utils.Shuffle(propositionChoices.Vocabs)
	utils.Shuffle(propositionChoices.Dependencies)

	return propositionChoices
}
