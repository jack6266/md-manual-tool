package input

import (
	"bufio"
	"fmt"
	"md-manual-tool/pkg/constants"
	"md-manual-tool/pkg/utils"
	"strings"
)

// Collector 输入收集器
type Collector struct {
	reader       *bufio.Reader
	versionUtils *utils.VersionUtils
}

// NewCollector 创建新的输入收集器
func NewCollector(reader *bufio.Reader) *Collector {
	return &Collector{
		reader:       reader,
		versionUtils: utils.NewVersionUtils(),
	}
}

// InputData 输入数据结构
type InputData struct {
	TemplatePath string
	ConfigPath   string
	Version      string
}

// CollectAll 收集所有输入
func (c *Collector) CollectAll() (*InputData, error) {
	data := &InputData{}

	// 收集模板文件路径
	if err := c.collectTemplatePath(data); err != nil {
		return nil, err
	}

	// 收集配置文件路径
	if err := c.collectConfigPath(data); err != nil {
		return nil, err
	}

	// 显示检测到的版本号
	c.showDetectedVersion(data.TemplatePath)

	// 收集版本号
	if err := c.collectVersion(data); err != nil {
		return nil, err
	}

	return data, nil
}

// collectTemplatePath 收集模板文件路径
func (c *Collector) collectTemplatePath(data *InputData) error {
	fmt.Print(constants.PromptTemplatePath)
	templatePath, err := c.reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf(constants.ErrReadTemplatePath, err)
	}
	data.TemplatePath = strings.TrimSpace(templatePath)
	return nil
}

// collectConfigPath 收集配置文件路径
func (c *Collector) collectConfigPath(data *InputData) error {
	fmt.Print(constants.PromptConfigPath)
	configPath, err := c.reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf(constants.ErrReadConfigPath, err)
	}
	data.ConfigPath = strings.TrimSpace(configPath)

	// 如果用户没有输入配置文件路径，则默认使用当前目录的config.yaml
	if data.ConfigPath == "" {
		data.ConfigPath = constants.DefaultConfigFile
		fmt.Printf(constants.MsgDefaultConfigUsed, data.ConfigPath)
	}
	return nil
}

// collectVersion 收集版本号
func (c *Collector) collectVersion(data *InputData) error {
	fmt.Print(constants.PromptVersion)
	version, err := c.reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf(constants.ErrReadVersion, err)
	}
	data.Version = strings.TrimSpace(version)

	// 验证版本号格式
	if !c.versionUtils.IsValidVersionFormat(data.Version) {
		return fmt.Errorf("版本号格式无效，请使用 x.y.z 格式（如 1.0.1）")
	}

	return nil
}

// showDetectedVersion 显示检测到的版本号
func (c *Collector) showDetectedVersion(templatePath string) {
	oldVersion := c.versionUtils.ExtractVersionFromFilename(templatePath)
	if oldVersion != "" {
		fmt.Printf(constants.MsgVersionDetected, oldVersion)
	}
}
