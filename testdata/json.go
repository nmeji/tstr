package testdata

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type jsonUnmarshaller struct {
	source string
}

func (u jsonUnmarshaller) unmarshal(m interface{}) error {
	f, err := os.Open(u.source)
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, m)
}

func newJSONUnmarshaller(sourceFilepath string) unmarshaller {
	return &jsonUnmarshaller{sourceFilepath}
}
