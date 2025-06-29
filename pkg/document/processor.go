package document

import (
	"fmt"
	"md-manual-tool/pkg/config"
	"md-manual-tool/pkg/constants"
	"md-manual-tool/pkg/processor"
)

// Processor 文档处理器
type Processor struct{}

// NewProcessor 创建新的文档处理器
func NewProcessor() *Processor {
	return &Processor{}
}

// ProcessDocument 处理文档
func (p *Processor) ProcessDocument(configData *config.ConfigData) error {
	// 创建处理器
	proc := processor.NewProcessor(configData.Config)

	// 处理整个流程
	if err := proc.Process(configData.TemplatePath, configData.OutputPath); err != nil {
		return fmt.Errorf(constants.ErrProcessFailed, err)
	}

	return nil
}

// ProcessWithConfig 使用配置处理文档
func (p *Processor) ProcessWithConfig(cfg *config.Config, templatePath, outputPath string) error {
	// 创建处理器
	proc := processor.NewProcessor(cfg)

	// 处理整个流程
	if err := proc.Process(templatePath, outputPath); err != nil {
		return fmt.Errorf(constants.ErrProcessFailed, err)
	}

	return nil
}
