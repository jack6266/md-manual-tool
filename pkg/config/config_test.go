package config

import (
	"os"
	"testing"
)

func TestReadConfig(t *testing.T) {
	// 创建临时配置文件
	content := `title: 测试项目
description: 这是一个测试描述
author: 测试作者
email: test@example.com`

	tmpFile, err := os.CreateTemp("", "test_config_*.yaml")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("写入临时文件失败: %v", err)
	}
	tmpFile.Close()

	// 测试读取配置
	config, err := ReadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("读取配置失败: %v", err)
	}

	// 验证配置内容
	expected := map[string]string{
		"title":       "测试项目",
		"description": "这是一个测试描述",
		"author":      "测试作者",
		"email":       "test@example.com",
	}

	for key, expectedValue := range expected {
		if actualValue, exists := config.Variables[key]; !exists {
			t.Errorf("缺少配置项: %s", key)
		} else if actualValue != expectedValue {
			t.Errorf("配置项 %s 的值不匹配，期望: %s, 实际: %s", key, expectedValue, actualValue)
		}
	}
}
