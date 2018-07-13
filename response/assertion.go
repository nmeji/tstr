package response

import (
	"fmt"
	"net/http"
	"reflect"
)

type Assert func() error

func assertEqual(expected, actual interface{}) error {
	if actual != expected {
		return fmt.Errorf("expected %v, got %v", expected, actual)
	}
	return nil
}

func assertDeepEqual(expected, actual interface{}) error {
	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected %v, got %v", expected, actual)
	}
	return nil
}

func NewChecker(resp *http.Response) *Checker {
	checker := &Checker{Response: resp}
	checker.ExpectBody = &BodyChecker{Checker: checker}
	return checker
}
