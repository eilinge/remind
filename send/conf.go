package send

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

func LoadConf(path string) (*Conf, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c Conf
	if err := yaml.Unmarshal(d, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

type Conf struct {
	Smtp       *SmtpConf
	HttpListen string `yaml:"http_listen"`
}

type SmtpConf struct {
	Addr     string `yaml:"addr"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
}
