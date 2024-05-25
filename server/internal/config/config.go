package config

type Config struct {
	Server ServerConfig `yaml:"server"`
	DB     *DBConfig    `yaml:"db"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
}
