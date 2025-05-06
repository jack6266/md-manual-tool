package main

import (
	"flag"
	"fmt"
	"md-manual-tool/config"
	"md-manual-tool/template"
	"md-manual-tool/utils"
)

func main() {
	// 定义命令行参数
	templatePath := flag.String("template", "templates/template.md", "模板文件路径")
	configPath := flag.String("config", "configs/config.yaml", "配置文件路径")
	outputPath := flag.String("output", "output/result.md", "输出文件路径")
	flag.Parse()

	// 读取配置文件
	config, err := config.ReadConfig(*configPath)
	if err != nil {
		fmt.Printf("读取配置文件失败: %v\n", err)
		return
	}

	// 渲染模板
	result, err := template.Render(*templatePath, config.Variables)
	if err != nil {
		fmt.Printf("渲染模板失败: %v\n", err)
		return
	}

	// 确保输出目录存在
	err = utils.EnsureDir(*outputPath)
	if err != nil {
		fmt.Printf("创建输出目录失败: %v\n", err)
		return
	}

	// 提取并复制图片，并更新Markdown内容
	imagePaths := utils.ExtractImages(string(result))
	if len(imagePaths) > 0 {
		updatedContent, err := utils.CopyImages(*outputPath, imagePaths, string(result))
		if err != nil {
			fmt.Printf("复制图片失败: %v\n", err)
			return
		}
		result = []byte(updatedContent)
		fmt.Printf("成功复制 %d 张图片\n", len(imagePaths))
	}

	// 写入结果文件
	err = utils.WriteFile(*outputPath, result)
	if err != nil {
		fmt.Printf("写入结果文件失败: %v\n", err)
		return
	}

	fmt.Printf("文件生成成功！输出路径：%s\n", *outputPath)
}
