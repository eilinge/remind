package send

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

func LoadConfUsers(path string) (UserSource, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var us ConfUserSource
	if err := yaml.Unmarshal(d, &us); err != nil {
		return nil, err
	}

	return us, nil
}

type User struct {
	Name   string `json:"name" yaml:"name"`
	Email  string `json:"email" yaml:"email"`
	Wechat string `json:"wechat" yaml:"wechat"`
}

type UserSource interface {
	Get(n *Notify) ([]*User, error)
}

type ConfUserSource map[string][]*User

func (c ConfUserSource) Get(n *Notify) ([]*User, error) {
	us, ok := c[n.Name]
	if !ok {
		return nil, errors.New("not found")
	}

	return us, nil
}
