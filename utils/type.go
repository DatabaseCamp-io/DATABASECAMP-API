package utils

import "encoding/json"

type typeUtils struct{}

type IType interface {
	StructToStruct(x1 interface{}, x2 interface{}) error
	StructToMap(obj interface{}) (newMap map[string]interface{}, err error)
}

func NewType() typeUtils {
	return typeUtils{}
}

func (t typeUtils) StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj)

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &newMap)
	return
}

func (t typeUtils) StructToStruct(x1 interface{}, x2 interface{}) error {
	temp, err := t.StructToMap(x1)
	if err != nil {
		return err
	}

	jsonBody, err := json.Marshal(temp)
	if err != nil {

		return err
	}

	if err := json.Unmarshal(jsonBody, x2); err != nil {

		return err
	}
	return nil
}
