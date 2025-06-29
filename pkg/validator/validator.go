package validator

import (
	"fmt"
	"md-manual-tool/pkg/constants"
	"os"
	"regexp"
)

// Validator 验证器
type Validator struct{}

// NewValidator 创建新的验证器
func NewValidator() *Validator {
	return &Validator{}
}

// ValidationResult 验证结果
type ValidationResult struct {
	IsValid bool
	Errors  []string
}

// NewValidationResult 创建新的验证结果
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		IsValid: true,
		Errors:  make([]string, 0),
	}
}

// ValidateInputs 验证所有输入
func (v *Validator) ValidateInputs(templatePath, configPath string) *ValidationResult {
	result := NewValidationResult()

	// 验证模板文件
	if err := v.validateFile(templatePath, constants.FileTypeTemplate); err != nil {
		result.IsValid = false
		result.Errors = append(result.Errors, err.Error())
	}

	// 验证配置文件
	if err := v.validateFile(configPath, constants.FileTypeConfig); err != nil {
		result.IsValid = false
		result.Errors = append(result.Errors, err.Error())
	}

	return result
}

// validateFile 验证文件是否存在
func (v *Validator) validateFile(filePath, fileType string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf(constants.ErrFileNotExist, fileType, filePath, v.getCurrentDir())
	}
	return nil
}

// getCurrentDir 获取当前工作目录
func (v *Validator) getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "无法获取当前目录"
	}
	return dir
}

// ValidateVersionFormat 验证版本号格式
func (v *Validator) ValidateVersionFormat(version string) error {
	if version == "" {
		return nil // 空版本号是允许的
	}

	// 验证版本号格式：x.y.z
	re := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	if !re.MatchString(version) {
		return fmt.Errorf("版本号格式无效，请使用 x.y.z 格式（如 1.0.1）")
	}

	return nil
}
