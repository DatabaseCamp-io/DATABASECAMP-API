package utils

// util.type.go
/**
 * 	This file is a part of utilities, used to convert type
 */

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

/**
 * 	This class convert any type
 */
type typeUtils struct{}

/**
 * Constructor creates a new typeUtils instance
 *
 * @return 	instance of typeUtils
 */
func NewType() typeUtils {
	return typeUtils{}
}

/**
 * Convert struct to map
 *
 * @param 	obj 	StructToMap to be converted to map
 *
 * @return 	map from obj converted
 * @return 	the error of converting
 */
func (t typeUtils) StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj)

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &newMap)
	return
}

/**
 * Convert struct to struct
 *
 * @param 	x1 	Original struct to be converted
 * @param 	x2 	Output struct
 *
 * @return 	the error of converting
 */
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

/**
 * Convert any type to duration
 *
 * @param 	item 	Item to be converted
 *
 * @return 	duration by converted item
 */
func (t typeUtils) ParseDuration(item interface{}) time.Duration {
	intType := t.ParseInt(item)
	return time.Duration(intType)
}

/**
 * Convert any type to integer value
 *
 * @param 	item 	Item to be converted
 *
 * @return 	integer value by converted item
 */
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

/**
 * Convert any type to string value
 *
 * @param 	item 	Item to be converted
 *
 * @return 	string value by converted item
 */
func (t typeUtils) ParseString(item interface{}) string {
	return fmt.Sprintf("%v", item)
}
