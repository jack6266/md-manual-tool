package config

import (
	"io/ioutil"
	"strings"
)

// Config 配置结构体
type Config struct {
	Variables map[string]string
}

// ReadConfig 读取配置文件
func ReadConfig(configPath string) (*Config, error) {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := &Config{
		Variables: make(map[string]string),
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			config.Variables[key] = value
		}
	}

	return config, nil
}
