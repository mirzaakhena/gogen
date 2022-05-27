package utils

import "encoding/json"

// MustJSON is converter from interface{} to string
// Warning! this function will always assume the conversion is success
// if you are not sure the conversion is always succeed then use ToJSON
func MustJSON(obj interface{}) string {
	bytes, _ := json.Marshal(obj)
	return string(bytes)
}
