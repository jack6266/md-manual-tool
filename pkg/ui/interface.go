package ui

import (
	"fmt"
	"md-manual-tool/pkg/constants"
)

// Interface UI交互接口
type Interface struct{}

// NewInterface 创建新的UI交互接口
func NewInterface() *Interface {
	return &Interface{}
}

// ShowSuccess 显示成功信息
func (ui *Interface) ShowSuccess(outputPath string) {
	fmt.Printf(constants.MsgFileGenerated, outputPath)
}

// ShowError 显示错误信息
func (ui *Interface) ShowError(message string) {
	fmt.Printf("错误：%s\n", message)
}

// ShowValidationErrors 显示验证错误
func (ui *Interface) ShowValidationErrors(errors []string) {
	fmt.Println("验证失败：")
	for _, err := range errors {
		fmt.Printf("  - %s\n", err)
	}
}

// ShowProgress 显示进度信息
func (ui *Interface) ShowProgress(message string) {
	fmt.Printf("正在处理：%s\n", message)
}

// ShowInfo 显示信息
func (ui *Interface) ShowInfo(message string) {
	fmt.Println(message)
}

// ShowInfoWithFormat 显示格式化信息
func (ui *Interface) ShowInfoWithFormat(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
