package conf

import (
	"io/ioutil"
	"money/pkg/mongodb"

	"money/pkg/log"

	"gopkg.in/yaml.v2"
)

var (
	ApolloPath string // refer by main
	ConfigPath string // refer by main

	global *Config

	//	errInvalidParam = errors.New("err invalid parmas")
)

type Config struct {
	Mongodb mongodb.Config `yaml:"mongodb"`
	Log     log.Config
	///Redis   redisdb.Config `yaml:"redis"`
}

func Init() *Config {
	if ApolloPath == "" && ConfigPath == "" {

		log.Fatal("apoll or local config are null")
	}

	err := LoadGlobal(ConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	return global
}
func LoadGlobal(fpath string) error {
	cfg, err := parse(fpath)
	if err != nil {
		return err
	}

	global = cfg
	return nil
}

func LoadGlobalContent(content string) error {
	c, err := parseContent(content)
	if err != nil {
		return err
	}
	global = c
	return nil
}

func parseContent(content string) (*Config, error) {
	var cfg Config
	err := yaml.Unmarshal([]byte(content), &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func GetConfig() *Config {
	if global == nil {
		return &Config{}
	}
	return global
}

func parse(fpath string) (*Config, error) {
	var cfg Config
	bs, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bs, &cfg)
	return &cfg, err
}
