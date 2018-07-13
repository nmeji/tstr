package response

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/oliveagle/jsonpath"
)

type HttpStatus struct {
	Text string
	Code int
}

type Checker struct {
	Response   *http.Response
	ExpectBody *BodyChecker
	asserts    []Assert
}

type BodyChecker struct {
	*Checker
	bodyCache []byte
}

func (b *BodyChecker) readBody() ([]byte, error) {
	if b.bodyCache == nil {
		defer b.Response.Body.Close()
		content, err := ioutil.ReadAll(b.Response.Body)
		if err != nil {
			return nil, err
		}
		b.bodyCache = content
	}
	return b.bodyCache, nil
}

func (b *BodyChecker) assertJson(expected interface{}, actualRaw []byte) error {
	t := reflect.TypeOf(expected)
	if t.Kind() == reflect.String {
		expectedJSONString, _ := expected.(string)
		return assertEqual(expectedJSONString, string(actualRaw))
	}
	actualPtr := reflect.New(t).Interface()
	err := json.Unmarshal(actualRaw, actualPtr)
	if err != nil {
		return err
	}
	actual := reflect.ValueOf(actualPtr).Elem().Interface()
	return assertDeepEqual(expected, actual)
}

func (b *BodyChecker) assertJsonWithPath(expectedInterface interface{}, content []byte, path string) error {
	var data interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return err
	}
	actualInterface, err := jsonpath.JsonPathLookup(data, path)
	if err != nil {
		return err
	}
	typeActual := reflect.TypeOf(actualInterface).Kind().String()
	switch expected := expectedInterface.(type) {
	case string:
		actual, ok := actualInterface.(string)
		if !ok {
			return fmt.Errorf("expected string, got %s", typeActual)
		}
		return assertEqual(expected, actual)
	case int, int8, int16, int32, int64:
		actual, ok := actualInterface.(int)
		if !ok {
			return fmt.Errorf("expected %T, got %s", expected, typeActual)
		}
		return assertEqual(expected, actual)
	case bool:
		actual, ok := actualInterface.(bool)
		if !ok {
			return fmt.Errorf("expected bool, got %s", typeActual)
		}
		return assertEqual(expected, actual)
	case float32, float64:
		actual, ok := actualInterface.(float64)
		if !ok {
			return fmt.Errorf("expected float64, got %s", typeActual)
		}
		return assertEqual(expected, actual)
	}
	actualJSON, err := json.Marshal(actualInterface)
	if err != nil {
		return err
	}
	return b.assertJson(expectedInterface, actualJSON)
}

func (b *BodyChecker) ToStringEqual(body string) *Checker {
	assert := func() error {
		content, err := b.readBody()
		if err != nil {
			return err
		}
		return assertEqual(body, string(content))
	}
	b.asserts = append(b.asserts, assert)
	return b.Checker
}

func (b *BodyChecker) ToJsonEqual(json interface{}) *Checker {
	assert := func() error {
		content, err := b.readBody()
		if err != nil {
			return err
		}
		return b.assertJson(json, content)
	}
	b.asserts = append(b.asserts, assert)
	return b.Checker
}

func (b *BodyChecker) ToHaveInJson(path string, json interface{}) *Checker {
	assert := func() error {
		content, err := b.readBody()
		if err != nil {
			return err
		}
		return b.assertJsonWithPath(json, content, path)
	}
	b.asserts = append(b.asserts, assert)
	return b.Checker
}

func (c *Checker) ExpectStatus(status int) *Checker {
	assert := func() error {
		return assertEqual(status, c.Response.StatusCode)
	}
	c.asserts = append(c.asserts, assert)
	return c
}

func (c Checker) MakeAssertion(t *testing.T) {
	for _, assert := range c.asserts {
		if err := assert(); err != nil {
			t.Errorf("%v", err)
		}
	}
}
