package environment

import (
	"github.com/joho/godotenv"
)

type dotEnv struct{}

func New() *dotEnv {
	return &dotEnv{}
}

func (env dotEnv) Load(path string) error {
	return godotenv.Load(path)
}
