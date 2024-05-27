package config

type Config struct {
	Server ServerConfig `yaml:"server"`
	DB     *DBConfig    `yaml:"db"`
	JWT    *JWTConfig   `yaml:"jwt"`
	Queue  *QueueConfig `yaml:"queue"`
	Email  *EmailConfig `yaml:"email"`
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

type JWTConfig struct {
	Secret string `yaml:"secret"`
	TTL    int    `yaml:"ttl"`
	Unit   string `yaml:"unit"`
}

type QueueConfig struct {
	Protocol string `yaml:"protocol"`
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
}

type EmailConfig struct {
	Domain string `yaml:"domain"`
	APIKey string `yaml:"apikey"`
}
