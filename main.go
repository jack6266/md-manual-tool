package main

import (
	"bufio"
	"fmt"
	"md-manual-tool/pkg/config"
	"md-manual-tool/pkg/processor"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("请输入模板文件路径（如 templates/template.md）：")
	templatePath, _ := reader.ReadString('\n')
	templatePath = strings.TrimSpace(templatePath)

	fmt.Print("请输入配置文件路径（如 D:/config.yaml，直接回车则使用工具的当前目录的config.yaml）：")
	configPath, _ := reader.ReadString('\n')
	configPath = strings.TrimSpace(configPath)

	// 如果用户没有输入配置文件路径，则默认使用当前目录的config.yaml
	if configPath == "" {
		configPath = "config.yaml"
		fmt.Printf("使用默认配置文件：%s\n", configPath)
	}

	fmt.Print("请输入输出文件路径（如 output/result.md）：")
	outputPath, _ := reader.ReadString('\n')
	outputPath = strings.TrimSpace(outputPath)

	// 检查文件是否存在
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		fmt.Printf("错误：模板文件不存在: %s\n", templatePath)
		fmt.Printf("当前工作目录: %s\n", getCurrentDir())
		return
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("错误：配置文件不存在: %s\n", configPath)
		fmt.Printf("当前工作目录: %s\n", getCurrentDir())
		return
	}

	// 读取配置文件
	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		fmt.Printf("读取配置文件失败: %v\n", err)
		return
	}

	// 创建处理器
	proc := processor.NewProcessor(cfg)

	// 处理整个流程
	if err := proc.Process(templatePath, outputPath); err != nil {
		fmt.Printf("处理失败: %v\n", err)
		return
	}

	fmt.Printf("文件生成成功！输出路径：%s\n", outputPath)
}

// getCurrentDir 获取当前工作目录
func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "无法获取当前目录"
	}
	return dir
}
