package activity

import (
	"database-camp/internal/errs"
	"database-camp/internal/utils"
)

func FormatAnswer(answer interface{}, activityTypeID int) (Answer, error) {
	switch activityTypeID {
	case 1:
		var matchingChoiceAnswer MatchingChoiceAnswer
		err := utils.StructToStruct(answer, &matchingChoiceAnswer)
		return matchingChoiceAnswer, err
	case 2:
		var multipleChoiceAnswer MultipleChoiceAnswer
		err := utils.StructToStruct(answer, &multipleChoiceAnswer)
		return multipleChoiceAnswer, err
	case 3:
		var completionChoiceAnswer CompletionChoiceAnswer
		err := utils.StructToStruct(answer, &completionChoiceAnswer)
		return completionChoiceAnswer, err
	case 4:
		var vocabGroupChoiceAnswer VocabGroupChoiceAnswer
		err := utils.StructToStruct(answer, &vocabGroupChoiceAnswer)
		return vocabGroupChoiceAnswer, err
	case 5:
		var dependencyChoiceAnswer DependencyChoiceAnswer
		err := utils.StructToStruct(answer, &dependencyChoiceAnswer)
		return dependencyChoiceAnswer, err
	default:
		return nil, errs.ErrActivityTypeInvalid
	}
}

type Answer interface {
	IsCorrect(choices Choices) (bool, error)
}

type MatchingItem struct {
	Item1 string `json:"item1"`
	Item2 string `json:"item2"`
}

type MatchingChoiceAnswer []MatchingItem

func (answer MatchingChoiceAnswer) IsCorrect(choices Choices) (bool, error) {
	matchingChoices, ok := choices.(MatchingChoices)
	if !ok {
		return false, errs.ErrAnswerInvalid
	}

	Item1Item2Map := map[string]string{}
	for _, correct := range matchingChoices {
		Item1Item2Map[correct.PairItem1] = correct.PairItem2
	}

	for _, item := range answer {
		if Item1Item2Map[item.Item1] != item.Item2 && Item1Item2Map[item.Item2] != item.Item1 {
			return false, nil
		}
	}

	return true, nil
}

type MultipleChoiceAnswer []int

func (answer MultipleChoiceAnswer) IsCorrect(choices Choices) (bool, error) {
	multipleChoices, ok := choices.(MultipleChoices)
	if !ok {
		return false, errs.ErrAnswerInvalid
	}

	countCorrect := 0

	solution := map[int]bool{}
	for _, choice := range multipleChoices {
		solution[choice.ID] = choice.IsCorrect

		if choice.IsCorrect {
			countCorrect++
		}
	}

	if countCorrect != len(answer) {
		return false, nil
	}

	for _, v := range answer {

		if _, ok := solution[v]; !ok {
			return false, nil
		}

		if !solution[v] {
			return false, nil
		}

	}

	return true, nil
}

type completionItem struct {
	ID      *int    `json:"completion_choice_id"`
	Content *string `json:"content"`
}

type CompletionChoiceAnswer []completionItem

func (answer CompletionChoiceAnswer) IsCorrect(choices Choices) (bool, error) {
	completionChoices, ok := choices.(CompletionChoices)
	if !ok {
		return false, errs.ErrAnswerInvalid
	}

	for _, choice := range completionChoices {
		for _, v := range answer {
			if (choice.ID == *v.ID) && (choice.Content != *v.Content) {
				return false, nil
			}
		}
	}

	return true, nil
}

type VocabGroupChoiceAnswer struct {
	Groups []VocabGroup `json:"groups"`
}

func (answer VocabGroupChoiceAnswer) IsCorrect(choices Choices) (bool, error) {

	vocabGroupChoice, ok := choices.(VocabGroupChoice)
	if !ok {
		return false, errs.ErrAnswerInvalid
	}

	if len(answer.Groups) != len(vocabGroupChoice.Groups) {
		return false, nil
	}

	solution := map[string]map[string]bool{}
	for _, choice := range vocabGroupChoice.Groups {
		if _, ok := solution[choice.GroupName]; !ok {
			solution[choice.GroupName] = map[string]bool{}
		}

		for _, vocab := range choice.Vocabs {
			solution[choice.GroupName][vocab] = true
		}
	}

	for _, group := range answer.Groups {

		if len(solution[group.GroupName]) != len(group.Vocabs) {
			return false, nil
		}

		for _, vocab := range group.Vocabs {
			if _, ok := solution[group.GroupName]; !ok {
				return false, nil
			}

			if !solution[group.GroupName][vocab] {
				return false, nil
			}
		}
	}

	return true, nil
}

type DependencyChoiceAnswer []Dependency

func (answer DependencyChoiceAnswer) IsCorrect(choices Choices) (bool, error) {
	choice, ok := choices.(DependencyChoice)
	if !ok {
		return false, errs.ErrAnswerInvalid
	}

	solution := map[string]map[string]bool{}
	for _, dependency := range choice.Dependencies {
		for _, determinant := range dependency.Determinants {
			if _, ok := solution[dependency.Dependent]; !ok {
				solution[dependency.Dependent] = map[string]bool{}
			}

			solution[dependency.Dependent][determinant.Value] = true
		}
	}

	for _, dependency := range answer {
		for _, determinant := range dependency.Determinants {
			if _, ok := solution[dependency.Dependent]; !ok {
				return false, nil
			}

			if !solution[dependency.Dependent][determinant.Value] {
				return false, nil
			}
		}
	}

	return true, nil
}

type ERChoiceAnswer struct {
	Tables        Tables        `json:"tables"`
	Relationships Relationships `json:"relationships"`
}

func (answer ERChoiceAnswer) IsCorrect(choice ERChoice) (bool, string) {

	if len(answer.Tables) < len(choice.Tables) {
		return false, RelationSuggestions[SUGGESTION_LESS_RELATION]
	}

	if len(answer.Tables) > len(choice.Tables) {
		return false, RelationSuggestions[SUGGESTION_MORE_RELATION]
	}

	if len(answer.Relationships) != len(choice.Relationships) {
		return false, RelationshipSuggestions[SUGGESTION_INCORRECT_NUMBER_RELATIONSHIP]
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

	if len(idMap) != len(answer.Tables) {
		return false, RelationSuggestions[SUGGESTION_DUPLICATION_RELATION]
	}

	for _, table := range answer.Tables {
		if _, ok := tableSolutionMap[table.Title]; !ok {
			return false, RelationSuggestions[SUGGESTION_INCORRECT_RELATION]
		} else {

			if len(table.Attributes) < len(tableSolutionMap[table.Title]) {
				return false, AttributeSuggestions[SUGGESTION_LESS_ATTRIBUTE]
			}

			if len(table.Attributes) > len(tableSolutionMap[table.Title]) {
				return false, AttributeSuggestions[SUGGESTION_MORE_ATTRIBUTE]
			}

			for _, attribute := range table.Attributes {
				if _, ok := tableSolutionMap[table.Title][attribute.Value]; !ok {
					return false, AttributeSuggestions[SUGGESTION_INCORRECT_ATTRIBUTE]
				} else {

					if tableSolutionMap[table.Title][attribute.Value].Key != nil && *tableSolutionMap[table.Title][attribute.Value].Key != *attribute.Key {
						return false, AttributeSuggestions[SUGGESTION_INCORRECT_KEY_ATTRIBUTE]
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
			return false, RelationshipSuggestions[SUGGESTION_INCORRECT_RELATIONSHIP]
		}
	}

	for _, s := range choice.Relationships {
		for _, a := range answer.Relationships {
			if s.Table1ID == idMap[a.Table1ID] && s.Table2ID == idMap[a.Table1ID] {
				return false, RelationshipSuggestions[SUGGESTION_INVALID_TYPE_RELATIONSHIP]
			}
		}
	}

	return true, ""
}
