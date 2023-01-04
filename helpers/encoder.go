package helpers

import (
	"encoding/json"
	"strconv"
)

func ConvertMap(oldMap map[string]interface{}) map[string]interface{} {
	s, _ := json.Marshal(oldMap)

	var raw map[string]json.RawMessage
	err := json.Unmarshal(s, &raw)
	if err != nil {
		panic(err)
	}
	parsed := make(map[string]interface{}, len(raw))
	for key, val := range raw {
		s := string(val)
		i, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			parsed[key] = i
			continue
		}
		f, err := strconv.ParseFloat(s, 64)
		if err == nil {
			parsed[key] = f
			continue
		}
		var v interface{}
		err = json.Unmarshal(val, &v)
		if err == nil {
			parsed[key] = v
			continue
		}
		parsed[key] = val
	}

	return parsed
}
