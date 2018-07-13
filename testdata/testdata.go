package testdata

import (
	"fmt"
)

type unmarshaller interface {
	unmarshal(m interface{}) error
}

type TestData struct {
	u unmarshaller
}

func (t TestData) Unmarshal(m interface{}) error {
	return t.u.unmarshal(m)
}

func ext(filePath string, length int) string {
	return filePath[len(filePath)-length:]
}

func New(sourceFilepath string) (*TestData, error) {
	switch {
	case ext(sourceFilepath, 5) == ".json":
		return &TestData{u: newJSONUnmarshaller(sourceFilepath)}, nil
	case ext(sourceFilepath, 4) == ".csv":
		return &TestData{u: newCSVUnmarshaller(sourceFilepath)}, nil
	case ext(sourceFilepath, 5) == ".toml":
		return &TestData{u: newTOMLUnmarshaller(sourceFilepath)}, nil
	}
	return nil, fmt.Errorf("Unknown format: %s", sourceFilepath)
}
