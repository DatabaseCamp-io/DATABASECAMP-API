package controller

type TodoResponse struct {
	Name string `json:"name"`
}

type ITodoController interface {
	GetAll() ([]TodoResponse, error)
}
