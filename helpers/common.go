package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gateway-api/helpers/utils"
)

func GoEnv() string {
	env := os.Getenv("GO_ENV")
	if (utils.Strings{"production", "prod"}).Include(env) {
		return "production"
	}

	if (utils.Strings{"development", "dev"}).Include(env) {
		return "development"
	}

	if (utils.Strings{"local"}).Include(env) {
		return "local"
	}

	return "production"
}

func Dump(i interface{}) {
	fmt.Println(ToJSON(i, "\t"))
}

func ToJSON(i interface{}, indent string) string {
	s, _ := json.MarshalIndent(i, "", indent)
	return string(s)
}

func StringReplacer(val string, replacer map[string]string) string {
	for k, v := range replacer {
		val = strings.Replace(val, fmt.Sprintf("{%s}", k), v, -1)
	}
	return val
}
