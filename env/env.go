package env

import (
	"os"
	"strconv"
)

var strenvmap = map[string]string{}
var intenvmap = map[string]int{}
var boolenvmap = map[string]bool{}

// StrVal returns the string value of an environment variable
func StrVal(k string) string {
	if v, ok := strenvmap[k]; ok {
		return v
	}
	v, _ := getenv(k)
	strenvmap[k] = v
	return v
}

// IntVal returns the int value of an environment variable
func IntVal(k string) (int, error) {
	if v, ok := intenvmap[k]; ok {
		return v, nil
	}
	if v, ok := getenv(k); ok {
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}
		intenvmap[k] = i
		return i, nil
	}
	return 0, nil
}

// BoolVal returns the boolean value of an environment variable
func BoolVal(k string) (bool, error) {
	if v, ok := boolenvmap[k]; ok {
		return v, nil
	}
	if v, ok := getenv(k); ok {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return false, err
		}
		boolenvmap[k] = b
		return b, nil
	}
	return false, nil
}

func getenv(k string) (string, bool) {
	v := os.Getenv(k)
	return v, len(v) > 0
}
