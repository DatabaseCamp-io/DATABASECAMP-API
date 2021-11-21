package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type typeUtils struct{}

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

func (t typeUtils) ParseDuration(item interface{}) time.Duration {
	intType := t.ParseInt(item)
	return time.Duration(intType)
}

func (t typeUtils) ParseInt(item interface{}) int {
	_item := item
	switch item.(type) {
	case int:
		return _item.(int)
	case int64:
		return int(_item.(int64))
	case int32:
		return int(_item.(int32))
	case float64:
		return int(_item.(float64))
	case string:
		_int, _ := strconv.Atoi(_item.(string))
		return _int
	default:
		return 0
	}
}

func (t typeUtils) ParseString(item interface{}) string {
	return fmt.Sprintf("%v", item)
}
