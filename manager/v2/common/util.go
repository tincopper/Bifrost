package common

import (
	"strconv"
	"strings"
)

func ParseIntDef(v string, def int) int {
	v = strings.Trim(v, "")
	result, err := strconv.Atoi(v)
	if err != nil {
		return def
	} else {
		return result
	}
}
