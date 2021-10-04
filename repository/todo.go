package repository

type Todo struct {
	ID   int    `gorm:"primaryKey;column:todo_id" json:"todo_id"`
	Name string `gorm:"column:name" json:"name"`
}

type ITodoRepository interface {
	GetAll() ([]Todo, error)
}
