package encoding

import (
	"encoding/json"
	"schoolserver/common/utils/array"
)

func EncodeJson(data interface{}) ([]byte, error) {
	changeValue := array.ArrayToMap(data, "json")
	return json.Marshal(changeValue)
}

func DecodeJson(data []byte, value interface{}) error {
	var valueDynamic interface{}
	err := json.Unmarshal(data, &valueDynamic)
	if err != nil {
		return err
	}
	return array.MapToArray(valueDynamic, value, "json")
}
