package testdata

import "github.com/BurntSushi/toml"

type tomlUnmarshaller struct {
	source string
}

func (u *tomlUnmarshaller) unmarshal(m interface{}) error {
	_, err := toml.DecodeFile(u.source, m)
	return err
}

func newTOMLUnmarshaller(sourceFilepath string) unmarshaller {
	return &tomlUnmarshaller{sourceFilepath}
}
