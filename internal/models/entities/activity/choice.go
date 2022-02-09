package activity

import (
	"database-camp/internal/utils"
)

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

type VocabGroup struct {
	GroupName string   `json:"group_name"`
	Vocabs    []string `json:"vocabs"`
}

type VocabGroupChoice struct {
	Groups []VocabGroup `json:"groups"`
}

func (choice VocabGroupChoice) CreatePropositionChoices() interface{} {
	groups := make([]string, 0)
	vocabs := make([]string, 0)

	for _, v := range choice.Groups {
		groups = append(groups, v.GroupName)
		vocabs = append(vocabs, v.Vocabs...)

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

func (Determinant) TableName() string {
	return "Determinant"
}

type Dependency struct {
	Dependent    string        `gorm:"column:dependent" json:"dependent"`
	Fixed        bool          `gorm:"column:fixed" json:"-"`
	Determinants []Determinant `gorm:"foreignKey:determinant_id" json:"determinants"`
}

func (Dependency) TableName() string {
	return "Dependency"
}

type DependencyChoice struct {
	ID           int          `gorm:"column:dependency_choice_id"`
	Dependencies []Dependency `gorm:"foreignKey:dependency_id"`
}

func (DependencyChoice) TableName() string {
	return "DependencyChoice"
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

	vocabs := make([]string, 0)
	dependencies := make([]dependency, 0)

	for _, v := range choice.Dependencies {
		dependencyResult := dependency{
			DeterminantsCount: len(v.Determinants),
			Determinants:      make([]string, 0),
		}

		if v.Fixed {
			ch := v.Dependent
			dependencyResult.Dependent = &ch
		} else {
			vocabs = append(vocabs, v.Dependent)
		}

		for _, d := range v.Determinants {
			if d.Fixed {
				dependencyResult.Determinants = append(dependencyResult.Determinants, d.Value)
			} else {
				vocabs = append(vocabs, d.Value)
			}
		}

		dependencies = append(dependencies, dependencyResult)
	}

	utils.Shuffle(vocabs)
	utils.Shuffle(dependencies)

	propositionChoices := result{
		Vocabs:       vocabs,
		Dependencies: dependencies,
	}

	return propositionChoices
}

const (
	ER_CHOICE_FILL_TABLE = "FILL_TABLE"
	ER_CHOICE_DRAW       = "DRAW"
)

type ERChoice struct {
	Type          string        `gorm:"column:type" json:"-"`
	Tables        Tables        `json:"tables"`
	Relationships Relationships `json:"relationships"`
}

func (ERChoice) TableName() string { return "ERChoice" }

func (choice ERChoice) CreatePropositionChoices() interface{} {
	vocabs := make([]string, 0)

	type TableChoice struct {
		Title           *string    `json:"title"`
		AttributesCount *int       `json:"attributes_count"`
		Attributes      Attributes `json:"attributes"`
	}

	tablesChoice := make([]TableChoice, 0)

	for _, v := range choice.Tables {
		tableChoice := TableChoice{
			Attributes: make(Attributes, 0),
		}

		if choice.Type == ER_CHOICE_FILL_TABLE {
			count := len(v.Attributes)
			tableChoice.AttributesCount = &count
		}

		tableChoice.Attributes = make(Attributes, 0)

		if v.Fixed {
			tableChoice.Title = &v.Title
		} else {
			vocabs = append(vocabs, v.Title)
		}

		for _, a := range v.Attributes {
			if a.Fixed {
				tableChoice.Attributes = append(tableChoice.Attributes, a)
			} else {
				vocabs = append(vocabs, a.Value)
			}
		}

		utils.Shuffle(tableChoice.Attributes)

		tablesChoice = append(tablesChoice, tableChoice)
	}

	if choice.Type == ER_CHOICE_FILL_TABLE {

		utils.Shuffle(vocabs)
		utils.Shuffle(tablesChoice)

		return map[string]interface {
		}{
			"vocabs": vocabs,
			"tables": tablesChoice,
		}

	}

	relationships := choice.Relationships

	for i, v := range choice.Relationships {
		if !v.Fixed {
			relationships = append(relationships[:i], relationships[i+1:]...)
		}
	}

	utils.Shuffle(relationships)

	return map[string]interface {
	}{
		"tables":        tablesChoice,
		"relationships": relationships,
	}
}

type ERAnswer struct {
	Tables        Tables        `json:"tables"`
	Relationships Relationships `json:"relationships"`
}

func (choice ERAnswer) CreatePropositionChoices() interface{} {
	return map[string]interface{}{
		"tables":        choice.Tables,
		"relationships": choice.Relationships,
		// "problems":      PeerProblems,
	}
}
