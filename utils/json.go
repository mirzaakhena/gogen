package utils

import "encoding/json"

// MustJSON is converter from any to string
// Warning! this function will always assume the conversion is success
// if you are not sure the conversion is always succeed then use ToJSON
func MustJSON(obj any) string {
	bytes, _ := json.Marshal(obj)
	return string(bytes)
}
