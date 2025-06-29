package processor

import (
	"fmt"
	"md-manual-tool/pkg/config"
	"md-manual-tool/pkg/template"
	"md-manual-tool/pkg/utils"
)

// Processor 处理器结构体
type Processor struct {
	config *config.Config
}

// NewProcessor 创建新的处理器
func NewProcessor(config *config.Config) *Processor {
	return &Processor{
		config: config,
	}
}

// Process 处理整个流程
func (p *Processor) Process(templatePath, outputPath string) error {
	// 1. 渲染模板
	result, err := template.Render(templatePath, p.config.Variables)
	if err != nil {
		return fmt.Errorf("渲染模板失败: %v", err)
	}

	// 2. 确保输出目录存在
	if err := utils.EnsureDir(outputPath); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 3. 处理图片
	updatedContent, err := p.processImages(outputPath, string(result))
	if err != nil {
		return fmt.Errorf("处理图片失败: %v", err)
	}

	// 4. 写入结果文件
	if err := utils.WriteFile(outputPath, []byte(updatedContent)); err != nil {
		return fmt.Errorf("写入结果文件失败: %v", err)
	}

	return nil
}

// processImages 处理图片
func (p *Processor) processImages(outputPath, content string) (string, error) {
	// 提取图片路径
	imagePaths := utils.ExtractImages(content)
	if len(imagePaths) == 0 {
		return content, nil
	}

	// 复制图片并更新路径
	updatedContent, err := utils.CopyImages(outputPath, imagePaths, content)
	if err != nil {
		return content, err
	}

	fmt.Printf("成功复制 %d 张图片\n", len(imagePaths))
	return updatedContent, nil
}
