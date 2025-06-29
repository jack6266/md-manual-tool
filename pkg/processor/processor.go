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
	// 1. 读取原始模板内容
	templateContent, err := utils.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("读取模板文件失败: %v", err)
	}

	// 2. 从原始模板中提取图片路径（在版本号替换之前）
	imagePaths := utils.ExtractImages(string(templateContent))
	fmt.Println("从原始模板中提取到的图片路径：")
	for i, path := range imagePaths {
		fmt.Printf("  %d. %s\n", i+1, path)
	}

	// 3. 处理图片（复制到新目录）
	if len(imagePaths) > 0 {
		updatedContent, err := utils.CopyImagesFromTemplate(templatePath, outputPath, imagePaths, string(templateContent))
		if err != nil {
			return fmt.Errorf("处理图片失败: %v", err)
		}
		templateContent = []byte(updatedContent)
		fmt.Printf("成功复制 %d 张图片\n", len(imagePaths))
	}

	// 4. 确保输出目录存在
	if err := utils.EnsureDir(outputPath); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 5. 渲染模板（在图片处理之后）
	result, err := template.RenderWithContent(templatePath, string(templateContent), p.config.Variables)
	if err != nil {
		return fmt.Errorf("渲染模板失败: %v", err)
	}

	// 6. 写入结果文件
	if err := utils.WriteFile(outputPath, result); err != nil {
		return fmt.Errorf("写入结果文件失败: %v", err)
	}

	return nil
}

// processImages 处理图片（保留原有方法以兼容）
func (p *Processor) processImages(templatePath, outputPath, content string) (string, error) {
	// 提取图片路径
	imagePaths := utils.ExtractImages(content)
	// 新增：打印所有图片路径
	fmt.Println("模板中提取到的图片路径：")
	for i, path := range imagePaths {
		fmt.Printf("  %d. %s\n", i+1, path)
	}
	if len(imagePaths) == 0 {
		return content, nil
	}

	// 复制图片并更新路径
	updatedContent, err := utils.CopyImagesFromTemplate(templatePath, outputPath, imagePaths, content)
	if err != nil {
		return content, err
	}

	fmt.Printf("成功复制 %d 张图片\n", len(imagePaths))
	return updatedContent, nil
}
