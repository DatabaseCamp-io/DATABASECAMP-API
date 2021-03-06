package activity

const (
	ATTRIBUTE_KEY_PK = "PK"
	ATTRIBUTE_KEY_FK = "FK"
)

type Table struct {
	ID         string     `gorm:"column:table_id" json:"table_id"`
	Title      string     `gorm:"column:title" json:"title"`
	Fixed      bool       `gorm:"column:fixed" json:"-"`
	Attributes Attributes `gorm:"-" json:"attributes"`
}

type Tables []Table

type Attribute struct {
	ID      int     `gorm:"column:attribute_id" json:"attribute_id"`
	TableID string  `gorm:"column:table_id" json:"-"`
	Key     *string `gorm:"column:key" json:"key"`
	Value   string  `gorm:"column:value" json:"value"`
	Fixed   bool    `gorm:"column:fixed" json:"-"`
}

type Attributes []Attribute

const (
	RELATIONSHIP_MANY_TO_MANY = "MANY_TO_MANY"
	RELATIONSHIP_ONE_TO_MANY  = "ONE_TO_MANY"
	RELATIONSHIP_ONE_TO_ONE   = "ONE_TO_ONE"
)

type Relationship struct {
	ID               int    `gorm:"column:relationship_id" json:"relationship_id"`
	RelationshipType string `gorm:"column:relationship_type" json:"relationship_type"`
	Table1ID         string `gorm:"column:table1_id" json:"table1_id"`
	Table2ID         string `gorm:"column:table2_id" json:"table2_id"`
	Fixed            bool   `gorm:"column:fixed" json:"-"`
}

type Relationships []Relationship
