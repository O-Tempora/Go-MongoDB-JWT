package server

type Config struct {
	Port     string `yaml:"port"`
	DbHost   string `yaml:"dbhost"`
	DbPort   string `yaml:"dbport"`
	Database string `yaml:"database"`
}

func NewConfig() *Config {
	return &Config{
		Port: ":5005",
	}
}
