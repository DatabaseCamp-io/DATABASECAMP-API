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
	Suggestions map[int]string
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

var RelationSuggestions map[int]string = map[int]string{
	SUGGESTION_LESS_RELATION:        "จำนวนของ Relation น้อยเกินไป",
	SUGGESTION_MORE_RELATION:        "จำนวนของ Relation มากเกินไป",
	SUGGESTION_INCORRECT_RELATION:   "Relation ไม่สอดคล้องกับความต้องการของระบบ",
	SUGGESTION_DUPLICATION_RELATION: "มี Relation ซ้ำกัน",
}

var RelationshipSuggestions map[int]string = map[int]string{
	SUGGESTION_INCORRECT_NUMBER_RELATIONSHIP: "จำนวนของ Relationship ไม่ถูกต้อง",
	SUGGESTION_INCORRECT_RELATIONSHIP:        "Relationship ระหว่าง Relation ไม่ถูกต้อง",
	SUGGESTION_INVALID_TYPE_RELATIONSHIP:     "ประเภทของ Relationship ไม่ถูกต้อง",
}

var AttributeSuggestions map[int]string = map[int]string{
	SUGGESTION_LESS_ATTRIBUTE:          "จำนวนของ Attribute น้อยเกินไป",
	SUGGESTION_MORE_ATTRIBUTE:          "จำนวนของ Attribute มากเกินไป",
	SUGGESTION_INCORRECT_ATTRIBUTE:     "Attribute ไม่สอดคล้องกับความต้องการของระบบ",
	SUGGESTION_INCORRECT_KEY_ATTRIBUTE: "Key ของ Attribute ไม่ถูกต้อง",
}