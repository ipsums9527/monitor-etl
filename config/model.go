package config

type Config struct {
	Listen string         `yaml:"listen"`
	Port   int            `yaml:"port"`
	Api    map[string]any `yaml:"api"`
}
