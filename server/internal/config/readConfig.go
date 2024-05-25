package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func InitConfig(filePath string) (*Config, error) {
	ext := filepath.Ext(filePath)
	var cfg Config
	if ext == ".yaml" || ext == ".yml" {
		configFile, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("Config file read failure: %s", err.Error())
		}

		err = yaml.Unmarshal(configFile, &cfg)
		if err != nil {
			return nil, fmt.Errorf("invalid YAML config file: %s", err.Error())
		}

		return &cfg, nil
	} else {
		err := cfg.getDBEnvFromOS()
		if err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}

func (c *Config) getDBEnvFromOS() error {
	dbCfg := &DBConfig{}

	dbCfg.Host = os.Getenv("DB_HOST")
	if dbCfg.Host == "" {
		return fmt.Errorf("DB_HOST not set")
	}

	dbCfg.Port = os.Getenv("DB_PORT")
	if dbCfg.Port == "" {
		return fmt.Errorf("DB_PORT not set")
	}

	dbCfg.Username = os.Getenv("DB_USERNAME")
	if dbCfg.Username == "" {
		return fmt.Errorf("DB_USERNAME not set")
	}

	dbCfg.Password = os.Getenv("DB_PASSWORD")
	if dbCfg.Password == "" {
		return fmt.Errorf("DB_PASSWORD not set")
	}

	dbCfg.Database = os.Getenv("DB_DATABASE")
	if dbCfg.Database == "" {
		return fmt.Errorf("DB_DATABASE not set")
	}

	dbCfg.SSLMode = os.Getenv("DB_SSL_MODE")
	if dbCfg.SSLMode == "" {
		return fmt.Errorf("DB_SSL_MODE not set")
	}

	c.DB = dbCfg

	return nil
}
