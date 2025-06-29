package utils

import (
	"md-manual-tool/pkg/constants"
	"regexp"
	"strings"
)

// VersionUtils 版本号工具结构体
type VersionUtils struct{}

// NewVersionUtils 创建新的版本号工具实例
func NewVersionUtils() *VersionUtils {
	return &VersionUtils{}
}

// ExtractVersionFromFilename 从文件名中提取版本号
func (v *VersionUtils) ExtractVersionFromFilename(filename string) string {
	// 匹配文件名末尾的版本号格式：_x.y.z.md
	re := regexp.MustCompile(constants.VersionRegexPattern)
	matches := re.FindStringSubmatch(filename)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// GenerateOutputFilename 生成输出文件名，保持产品名称，替换版本号
func (v *VersionUtils) GenerateOutputFilename(templatePath, newVersion string) string {
	// 获取文件名（不含路径）
	filename := extractBaseFilename(templatePath)

	// 如果文件名包含版本号，替换为新版本号
	if oldVersion := v.ExtractVersionFromFilename(filename); oldVersion != "" && newVersion != "" {
		// 替换版本号
		newFilename := strings.Replace(filename, "_"+oldVersion+".md", "_"+newVersion+".md", 1)
		return newFilename
	}

	// 如果没有版本号或新版本号为空，使用默认名称
	if newVersion != "" {
		return newVersion + ".md"
	}
	return constants.DefaultOutputFile
}

// extractBaseFilename 提取基础文件名（不含路径）
func extractBaseFilename(filePath string) string {
	// 处理Windows和Unix路径分隔符
	parts := strings.Split(filePath, "/")
	parts = strings.Split(parts[len(parts)-1], "\\")
	return parts[len(parts)-1]
}

// IsValidVersionFormat 验证版本号格式
func (v *VersionUtils) IsValidVersionFormat(version string) bool {
	if version == "" {
		return true // 空版本号是允许的
	}
	// 验证版本号格式：x.y.z
	re := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	return re.MatchString(version)
}
