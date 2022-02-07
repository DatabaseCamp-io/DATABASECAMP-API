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

	for _, v := range answer {
		hasChoice := false

		for _, choice := range multipleChoices {
			if choice.ID == v {
				hasChoice = true

				if !choice.IsCorrect {
					return false, nil
				}

			}

			if !hasChoice {
				return false, nil
			}
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

type VocabGroup struct {
	GroupName string   `json:"group_name"`
	Vocabs    []string `json:"vocabs"`
}

type VocabGroupChoiceAnswer struct {
	Groups []VocabGroup `json:"groups"`
}

func (answer VocabGroupChoiceAnswer) IsCorrect(choices Choices) (bool, error) {
	vocabGroups, ok := choices.(VocalGroupChoices)
	if !ok {
		return false, errs.ErrAnswerInvalid
	}

	solution := make(map[string]map[string]bool, 0)
	for _, choice := range vocabGroups {
		if _, ok := solution[choice.GroupName]; !ok {
			solution[choice.GroupName] = make(map[string]bool, 0)
		}

		solution[choice.GroupName][choice.Vocab] = true
	}

	for _, group := range answer.Groups {
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

	solution := make(map[string]map[string]bool, 0)
	for _, dependency := range choice.Dependencies {
		for _, determinant := range dependency.Determinants {
			if _, ok := solution[dependency.Dependent]; !ok {
				solution[dependency.Dependent] = make(map[string]bool, 0)
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
