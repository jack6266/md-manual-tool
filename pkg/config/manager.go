package config

import (
	"fmt"
	"md-manual-tool/pkg/constants"
	"md-manual-tool/pkg/utils"
	"os"
	"path/filepath"
)

// Manager 配置管理器
type Manager struct {
	versionUtils *utils.VersionUtils
}

// NewManager 创建新的配置管理器
func NewManager() *Manager {
	return &Manager{
		versionUtils: utils.NewVersionUtils(),
	}
}

// ConfigData 配置数据结构
type ConfigData struct {
	Config       *Config
	OutputPath   string
	TemplatePath string
	Version      string
}

// LoadAndProcessConfig 加载并处理配置
func (m *Manager) LoadAndProcessConfig(configPath, templatePath, version string) (*ConfigData, error) {
	// 读取配置文件
	cfg, err := ReadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrReadConfig, err)
	}

	// 将版本参数添加到配置变量中
	if version != "" {
		cfg.Variables["version"] = version
		fmt.Printf(constants.MsgVersionAdded, version)
	}

	// 生成输出文件名
	outputFilename := m.versionUtils.GenerateOutputFilename(templatePath, version)

	// 创建输出目录路径：当前目录/output/
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("获取当前工作目录失败: %v", err)
	}
	outputPath := filepath.Join(currentDir, "output", outputFilename)

	return &ConfigData{
		Config:       cfg,
		OutputPath:   outputPath,
		TemplatePath: templatePath,
		Version:      version,
	}, nil
}

// AddVersionToConfig 将版本号添加到配置中
func (m *Manager) AddVersionToConfig(cfg *Config, version string) {
	if version != "" {
		cfg.Variables["version"] = version
	}
}

// GenerateOutputPath 生成输出路径
func (m *Manager) GenerateOutputPath(templatePath, version string) string {
	return m.versionUtils.GenerateOutputFilename(templatePath, version)
}
