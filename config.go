package synapse

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	PluginDir string          `yaml:"plugin_dir"`
	Matchers  []ConfigMatcher `yaml:"matchers"`
}

func (c *Config) LoadYAML(config string) error {
	f, err := ioutil.ReadFile(config)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(f, c)
	if err != nil {
		return err
	}
	return nil
}

type ConfigMatcher struct {
	Name       string         `yaml:"name"`
	Host       string         `yaml:"host"`
	Profilers  []ConfigPlugin `yaml:"profilers"`
	Associator ConfigPlugin   `yaml:"associator"`
	Searcher   ConfigPlugin   `yaml:"searcher"`
}

type ConfigPlugin struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}
