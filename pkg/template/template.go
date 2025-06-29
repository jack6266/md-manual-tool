package template

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"text/template"
)

// Render 渲染模板
func Render(templatePath string, variables map[string]string) ([]byte, error) {
	fmt.Printf("开始渲染模板: %s\n", templatePath)

	// 读取模板文件
	templateContent, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}
	fmt.Printf("模板文件大小: %d 字节\n", len(templateContent))

	// 从模板文件名中提取版本号
	oldVersion := extractVersionFromFilename(templatePath)
	if oldVersion != "" {
		fmt.Printf("从文件名提取的版本号: %s\n", oldVersion)
	}

	// 如果有新版本号且找到了原版本号，进行替换
	if newVersion, exists := variables["version"]; exists && oldVersion != "" {
		fmt.Printf("进行版本号替换: %s -> %s\n", oldVersion, newVersion)
		content := string(templateContent)
		// 替换内容中的版本号（排除图片路径）
		content = replaceVersionInContent(content, oldVersion, newVersion)
		templateContent = []byte(content)
		fmt.Printf("版本号替换完成\n")
	}

	// 创建模板
	tmpl, err := template.New("md").Parse(string(templateContent))
	if err != nil {
		return nil, err
	}

	// 渲染模板
	var result bytes.Buffer
	err = tmpl.Execute(&result, variables)
	if err != nil {
		return nil, err
	}

	fmt.Printf("模板渲染完成，结果大小: %d 字节\n", result.Len())
	return result.Bytes(), nil
}

// RenderWithContent 使用已读取的模板内容进行渲染
func RenderWithContent(templatePath string, templateContent string, variables map[string]string) ([]byte, error) {
	fmt.Printf("开始渲染模板内容: %s\n", templatePath)

	// 从模板文件名中提取版本号
	oldVersion := extractVersionFromFilename(templatePath)
	if oldVersion != "" {
		fmt.Printf("从文件名提取的版本号: %s\n", oldVersion)
	}

	// 如果有新版本号且找到了原版本号，进行替换
	if newVersion, exists := variables["version"]; exists && oldVersion != "" {
		fmt.Printf("进行版本号替换: %s -> %s\n", oldVersion, newVersion)
		// 替换内容中的版本号（排除图片路径）
		templateContent = replaceVersionInContent(templateContent, oldVersion, newVersion)
		fmt.Printf("版本号替换完成\n")
	}

	// 创建模板
	tmpl, err := template.New("md").Parse(templateContent)
	if err != nil {
		return nil, err
	}

	// 渲染模板
	var result bytes.Buffer
	err = tmpl.Execute(&result, variables)
	if err != nil {
		return nil, err
	}

	fmt.Printf("模板渲染完成，结果大小: %d 字节\n", result.Len())
	return result.Bytes(), nil
}

// extractVersionFromFilename 从文件名中提取版本号
func extractVersionFromFilename(filename string) string {
	// 匹配文件名末尾的版本号格式：_x.y.z.md
	re := regexp.MustCompile(`_(\d+\.\d+\.\d+)\.md$`)
	matches := re.FindStringSubmatch(filename)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// replaceVersionInContent 在内容中替换版本号（排除图片路径）
func replaceVersionInContent(content, oldVersion, newVersion string) string {
	// 先保护图片路径，避免被版本号替换影响
	protectedContent := protectImagePaths(content)

	// 替换各种可能的版本号格式
	replacements := []string{
		oldVersion,               // 直接替换
		"v" + oldVersion,         // v1.0.0 格式
		"版本 " + oldVersion,       // 版本 1.0.0 格式
		"Version " + oldVersion,  // Version 1.0.0 格式
		"版本号：" + oldVersion,      // 版本号：1.0.0 格式
		"Version: " + oldVersion, // Version: 1.0.0 格式
	}

	result := protectedContent
	for _, old := range replacements {
		new := strings.Replace(old, oldVersion, newVersion, -1)
		result = strings.Replace(result, new, old, -1) // 恢复被保护的图片路径
		result = strings.Replace(result, old, new, -1) // 进行版本号替换
	}

	// 恢复图片路径
	result = restoreImagePaths(result)

	return result
}

// protectImagePaths 保护图片路径，避免被版本号替换影响
func protectImagePaths(content string) string {
	// 匹配图片路径的正则表达式
	imagePattern := `!\[.*?\]\([^)]*\.(png|jpg|jpeg|gif|bmp|webp|svg|ico|tiff|tif)\)`
	re := regexp.MustCompile(imagePattern)

	// 用占位符替换图片路径
	imagePaths := make(map[string]string)
	counter := 0

	protectedContent := re.ReplaceAllStringFunc(content, func(match string) string {
		placeholder := fmt.Sprintf("__IMAGE_PATH_%d__", counter)
		imagePaths[placeholder] = match
		counter++
		return placeholder
	})

	// 将图片路径映射存储到全局变量中（这里简化处理，实际应该用更好的方式）
	globalImagePaths = imagePaths

	return protectedContent
}

// restoreImagePaths 恢复图片路径
func restoreImagePaths(content string) string {
	result := content
	for placeholder, imagePath := range globalImagePaths {
		result = strings.Replace(result, placeholder, imagePath, -1)
	}
	globalImagePaths = nil // 清理全局变量
	return result
}

// 全局变量存储图片路径映射（简化实现）
var globalImagePaths map[string]string
