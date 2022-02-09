package activity

const (
	SUGGESTION_LESS_RELATION = iota
	SUGGESTION_MORE_RELATION
	SUGGESTION_DUPLICATION_RELATION
	SUGGESTION_INCORRECT_RELATION
	SUGGESTION_INCORRECT_NUMBER_RELATIONSHIP
	SUGGESTION_INCORRECT_RELATIONSHIP
	SUGGESTION_INVALID_TYPE_RELATIONSHIP
	SUGGESTION_LESS_ATTRIBUTE
	SUGGESTION_MORE_ATTRIBUTE
	SUGGESTION_INCORRECT_ATTRIBUTE
	SUGGESTION_INCORRECT_KEY_ATTRIBUTE
)

type SuggestionGroup struct {
	Name        string
	Suggestions Suggestions
}

var SuggestionGroups = []SuggestionGroup{
	{
		Name:        "ด้าน Relation",
		Suggestions: RelationSuggestions,
	},
	{
		Name:        "ด้าน Relationship",
		Suggestions: RelationshipSuggestions,
	},
	{
		Name:        "ด้าน Attribute",
		Suggestions: AttributeSuggestions,
	},
}

type Suggestions map[int]string

func (s Suggestions) Strings() []string {
	suggestions := make([]string, 0)

	for _, v := range s {
		suggestions = append(suggestions, v)
	}

	return suggestions
}

var RelationSuggestions Suggestions = Suggestions{
	SUGGESTION_LESS_RELATION:        "จำนวนของ Relation น้อยเกินไป",
	SUGGESTION_MORE_RELATION:        "จำนวนของ Relation มากเกินไป",
	SUGGESTION_INCORRECT_RELATION:   "Relation ไม่สอดคล้องกับความต้องการของระบบ",
	SUGGESTION_DUPLICATION_RELATION: "มี Relation ซ้ำกัน",
}

var RelationshipSuggestions Suggestions = Suggestions{
	SUGGESTION_INCORRECT_NUMBER_RELATIONSHIP: "จำนวนของ Relationship ไม่ถูกต้อง",
	SUGGESTION_INCORRECT_RELATIONSHIP:        "Relationship ระหว่าง Relation ไม่ถูกต้อง",
	SUGGESTION_INVALID_TYPE_RELATIONSHIP:     "ประเภทของ Relationship ไม่ถูกต้อง",
}

var AttributeSuggestions Suggestions = Suggestions{
	SUGGESTION_LESS_ATTRIBUTE:          "จำนวนของ Attribute น้อยเกินไป",
	SUGGESTION_MORE_ATTRIBUTE:          "จำนวนของ Attribute มากเกินไป",
	SUGGESTION_INCORRECT_ATTRIBUTE:     "Attribute ไม่สอดคล้องกับความต้องการของระบบ",
	SUGGESTION_INCORRECT_KEY_ATTRIBUTE: "Key ของ Attribute ไม่ถูกต้อง",
}

type ProblemGroup struct {
	Name    string   `json:"name"`
	Choices []string `json:"choices"`
}

type PeerProblems struct {
	Groups []ProblemGroup `json:"groups"`
}

func GetPeerProblem() PeerProblems {
	peerProblems := PeerProblems{}

	for _, group := range SuggestionGroups {
		peerProblems.Groups = append(peerProblems.Groups, ProblemGroup{
			Name:    group.Name,
			Choices: group.Suggestions.Strings(),
		})
	}

	return peerProblems
}
