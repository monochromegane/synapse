package synapse

type Config struct {
	PluginDir string          `yaml:"plugin_dir"`
	Matchers  []ConfigMatcher `yaml:"matchers"`
}

type ConfigMatcher struct {
	Name       string         `yaml:"name"`
	Host       string         `yaml:"host"`
	Profiles   []ConfigPlugin `yaml:"profiles"`
	Associator ConfigPlugin   `yaml:"associator"`
	Searcher   ConfigPlugin   `yaml:"searcher"`
}

type ConfigPlugin struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}
