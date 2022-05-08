package activity

import (
	"database-camp/internal/logs"
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

type ProblemGroups []ProblemGroup

func (g ProblemGroups) Compare(list []string) bool {

	solution := map[string]bool{}
	for _, p := range g {
		for _, c := range p.Choices {
			solution[c] = true
		}
	}

	if len(solution) != len(list) {
		return false
	}

	for _, c := range list {
		if !solution[c] {
			return false
		}
	}

	return true
}

func (choice ERChoice) GetSuggestionsList(answer ERChoiceAnswer) ProblemGroups {
	problemMap := map[string][]string{
		SUGGESTION_RELATION_GROUP:     make([]string, 0),
		SUGGESTION_RELATIONSHIP_GROUP: make([]string, 0),
		SUGGESTION_ATTRIBUTE_GROUP:    make([]string, 0),
	}

	if len(answer.Tables) < len(choice.Tables) {
		problemMap[SUGGESTION_RELATION_GROUP] = append(problemMap[SUGGESTION_RELATION_GROUP], RelationSuggestions[SUGGESTION_LESS_RELATION])
	}

	if len(answer.Tables) > len(choice.Tables) {
		problemMap[SUGGESTION_RELATION_GROUP] = append(problemMap[SUGGESTION_RELATION_GROUP], RelationSuggestions[SUGGESTION_MORE_RELATION])
	}

	if len(answer.Relationships) != len(choice.Relationships) {
		problemMap[SUGGESTION_RELATIONSHIP_GROUP] = append(problemMap[SUGGESTION_RELATIONSHIP_GROUP], RelationshipSuggestions[SUGGESTION_INCORRECT_NUMBER_RELATIONSHIP])
	}

	tableSolutionMap := map[string]map[string]Attribute{}
	for _, table := range choice.Tables {
		tableSolutionMap[table.Title] = map[string]Attribute{}
		for _, attribute := range table.Attributes {
			tableSolutionMap[table.Title][attribute.Value] = attribute
		}

	}

	idMap := map[string]string{}
	for _, a := range answer.Tables {
		for _, s := range choice.Tables {
			if a.Title == s.Title {
				idMap[a.ID] = s.ID
			}
		}
	}

	logs.GetInstance().Info(idMap)

	if len(idMap) != len(choice.Tables) {
		problemMap[SUGGESTION_RELATION_GROUP] = append(problemMap[SUGGESTION_RELATION_GROUP], RelationSuggestions[SUGGESTION_INCORRECT_RELATION])
	}

	for _, table := range answer.Tables {
		if _, ok := tableSolutionMap[table.Title]; !ok {
			problemMap[SUGGESTION_RELATION_GROUP] = append(problemMap[SUGGESTION_RELATION_GROUP], RelationSuggestions[SUGGESTION_INCORRECT_RELATION])
		} else {

			if len(table.Attributes) < len(tableSolutionMap[table.Title]) {
				problemMap[SUGGESTION_ATTRIBUTE_GROUP] = append(problemMap[SUGGESTION_ATTRIBUTE_GROUP], AttributeSuggestions[SUGGESTION_LESS_ATTRIBUTE])
			}

			if len(table.Attributes) > len(tableSolutionMap[table.Title]) {
				problemMap[SUGGESTION_ATTRIBUTE_GROUP] = append(problemMap[SUGGESTION_ATTRIBUTE_GROUP], AttributeSuggestions[SUGGESTION_MORE_ATTRIBUTE])
			}

			for _, attribute := range table.Attributes {
				if _, ok := tableSolutionMap[table.Title][attribute.Value]; !ok {
					problemMap[SUGGESTION_ATTRIBUTE_GROUP] = append(problemMap[SUGGESTION_ATTRIBUTE_GROUP], AttributeSuggestions[SUGGESTION_INCORRECT_ATTRIBUTE])
				} else {

					if attribute.Key == nil && tableSolutionMap[table.Title][attribute.Value].Key != attribute.Key {
						problemMap[SUGGESTION_ATTRIBUTE_GROUP] = append(problemMap[SUGGESTION_ATTRIBUTE_GROUP], AttributeSuggestions[SUGGESTION_INCORRECT_KEY_ATTRIBUTE])
					} else if tableSolutionMap[table.Title][attribute.Value].Key != nil && *tableSolutionMap[table.Title][attribute.Value].Key != *attribute.Key {
						problemMap[SUGGESTION_ATTRIBUTE_GROUP] = append(problemMap[SUGGESTION_ATTRIBUTE_GROUP], AttributeSuggestions[SUGGESTION_INCORRECT_KEY_ATTRIBUTE])
					}

				}
			}

		}
	}

	relationshipMap := map[string]map[string]bool{}
	for _, r := range choice.Relationships {
		if _, ok := relationshipMap[r.Table1ID]; !ok {
			relationshipMap[r.Table1ID] = map[string]bool{}
		}

		relationshipMap[r.Table1ID][r.Table2ID] = true
	}

	for _, a := range answer.Relationships {
		if !relationshipMap[idMap[a.Table1ID]][idMap[a.Table2ID]] {
			problemMap[SUGGESTION_RELATIONSHIP_GROUP] = append(problemMap[SUGGESTION_RELATIONSHIP_GROUP], RelationshipSuggestions[SUGGESTION_INCORRECT_RELATIONSHIP])
		}
	}

	for _, s := range choice.Relationships {
		for _, a := range answer.Relationships {
			if s.Table1ID == idMap[a.Table1ID] && s.Table2ID == idMap[a.Table1ID] {
				problemMap[SUGGESTION_RELATIONSHIP_GROUP] = append(problemMap[SUGGESTION_RELATIONSHIP_GROUP], RelationshipSuggestions[SUGGESTION_INVALID_TYPE_RELATIONSHIP])
			}
		}
	}

	problemGroup := make([]ProblemGroup, 0)

	for i, v := range problemMap {
		problemGroup = append(problemGroup, ProblemGroup{
			Name:    i,
			Choices: v,
		})
	}

	logs.GetInstance().Info(problemGroup)

	return problemGroup

}

func (choice ERChoice) CreatePropositionChoices() interface{} {
	vocabs := make([]string, 0)

	type TableChoice struct {
		TableID         *string    `json:"table_id"`
		Title           *string    `json:"title"`
		AttributesCount *int       `json:"attributes_count"`
		Attributes      Attributes `json:"attributes"`
	}

	tablesChoice := make([]TableChoice, 0)

	for i, table := range choice.Tables {

		tableChoice := TableChoice{
			Attributes: []Attribute{},
		}

		if choice.Type == ER_CHOICE_FILL_TABLE {
			count := len(table.Attributes)
			tableChoice.AttributesCount = &count
		}

		if table.Fixed {
			tableChoice.TableID = &choice.Tables[i].ID
			tableChoice.Title = &choice.Tables[i].Title
		} else {
			vocabs = append(vocabs, table.Title)
		}

		for _, a := range table.Attributes {
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
	_relationships := make([]Relationship, 0)

	for i, v := range choice.Relationships {

		if v.Fixed {
			_relationships = append(_relationships, relationships[i])
		}

	}

	utils.Shuffle(_relationships)

	return map[string]interface {
	}{
		"tables":        tablesChoice,
		"relationships": _relationships,
	}
}

type ERAnswerTables struct {
	ERAnswerID int    `gorm:"column:er_answer_id" json:"-"`
	TableID    string `gorm:"column:table_id" json:"-"`
}

type ERAnswer struct {
	ID            int           `gorm:"column:er_answer_id" json:"er_answer_id"`
	UserID        int           `gorm:"column:user_id" json:"-"`
	Tables        Tables        `gorm:"-" json:"tables"`
	Relationships Relationships `gorm:"-" json:"relationships"`
}

func (choice ERAnswer) CreatePropositionChoices() interface{} {
	return map[string]interface{}{
		"er_answe_id":   choice.ID,
		"tables":        choice.Tables,
		"relationships": choice.Relationships,
		"problems":      GetPeerProblem(),
	}
}
