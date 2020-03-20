package send

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Content struct {
	Message struct {
		Cont  string `json:"content" yaml:"content"`
		Level string `json:"level" yaml:"level"`
	}
}

func LoadMessage(path string) (*Content, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c Content
	if err := yaml.Unmarshal(d, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
