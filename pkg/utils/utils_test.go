package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestCopyImagesFromTemplate(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "test_images")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建测试模板文件
	templatePath := filepath.Join(tempDir, "test_template.md")
	templateContent := `# 测试文档

![测试图片](./images/test.png)
![另一个图片](./images/another.jpg)
![相对路径图片](images/relative.jpg)
`
	err = os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("创建模板文件失败: %v", err)
	}

	// 创建图片目录和测试图片
	imagesDir := filepath.Join(tempDir, "images")
	err = os.MkdirAll(imagesDir, 0755)
	if err != nil {
		t.Fatalf("创建图片目录失败: %v", err)
	}

	// 创建测试图片文件
	testImage1 := filepath.Join(imagesDir, "test.png")
	err = os.WriteFile(testImage1, []byte("fake png content"), 0644)
	if err != nil {
		t.Fatalf("创建测试图片1失败: %v", err)
	}

	testImage2 := filepath.Join(imagesDir, "another.jpg")
	err = os.WriteFile(testImage2, []byte("fake jpg content"), 0644)
	if err != nil {
		t.Fatalf("创建测试图片2失败: %v", err)
	}

	testImage3 := filepath.Join(imagesDir, "relative.jpg")
	err = os.WriteFile(testImage3, []byte("fake relative content"), 0644)
	if err != nil {
		t.Fatalf("创建测试图片3失败: %v", err)
	}

	// 设置输出路径
	outputPath := filepath.Join(tempDir, "output.md")

	// 提取图片路径
	imagePaths := ExtractImages(templateContent)
	if len(imagePaths) != 3 {
		t.Errorf("期望提取3个图片路径，实际提取了%d个", len(imagePaths))
	}

	// 测试复制图片
	updatedContent, err := CopyImagesFromTemplate(templatePath, outputPath, imagePaths, templateContent)
	if err != nil {
		t.Fatalf("复制图片失败: %v", err)
	}

	// 检查输出内容是否包含正确的图片路径
	expectedPath1 := "./output.assets/test.png"
	expectedPath2 := "./output.assets/another.jpg"
	expectedPath3 := "./output.assets/relative.jpg"

	if !contains(updatedContent, expectedPath1) {
		t.Errorf("输出内容中未找到预期的图片路径: %s", expectedPath1)
	}

	if !contains(updatedContent, expectedPath2) {
		t.Errorf("输出内容中未找到预期的图片路径: %s", expectedPath2)
	}

	if !contains(updatedContent, expectedPath3) {
		t.Errorf("输出内容中未找到预期的图片路径: %s", expectedPath3)
	}

	// 检查图片文件是否被复制
	outputImagesDir := filepath.Join(tempDir, "output.assets")
	copiedImage1 := filepath.Join(outputImagesDir, "test.png")
	copiedImage2 := filepath.Join(outputImagesDir, "another.jpg")
	copiedImage3 := filepath.Join(outputImagesDir, "relative.jpg")

	if _, err := os.Stat(copiedImage1); os.IsNotExist(err) {
		t.Errorf("复制的图片文件不存在: %s", copiedImage1)
	}

	if _, err := os.Stat(copiedImage2); os.IsNotExist(err) {
		t.Errorf("复制的图片文件不存在: %s", copiedImage2)
	}

	if _, err := os.Stat(copiedImage3); os.IsNotExist(err) {
		t.Errorf("复制的图片文件不存在: %s", copiedImage3)
	}
}

func TestResolveImagePath(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "test_resolve")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建测试文件
	templatePath := filepath.Join(tempDir, "template.md")
	imagePath := filepath.Join(tempDir, "images", "test.png")

	// 创建图片目录和文件
	err = os.MkdirAll(filepath.Dir(imagePath), 0755)
	if err != nil {
		t.Fatalf("创建图片目录失败: %v", err)
	}

	err = os.WriteFile(imagePath, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("创建测试图片失败: %v", err)
	}

	// 测试相对路径解析
	resolvedPath, err := resolveImagePath("images/test.png", templatePath)
	if err != nil {
		t.Errorf("解析相对路径失败: %v", err)
	}

	expectedPath := filepath.Join(tempDir, "images", "test.png")
	if resolvedPath != expectedPath {
		t.Errorf("解析路径错误: 期望 %s, 实际 %s", expectedPath, resolvedPath)
	}
}

func TestExtractImages(t *testing.T) {
	// 测试内容
	content := `# 测试文档

![图片1](./images/test1.png)
![图片2](images/test2.jpg)
![图片3](../images/test3.gif)
![图片4](/absolute/path/image.png)
![图片5](images\\test5.png)
![图片6](images/test6.png?v=123)
![图片7](./易立德产品数据管理软件(eRDCloud-PDM)部署手册_3.2.1.assets/1.png)
`

	// 提取图片路径
	paths := ExtractImages(content)

	// 验证提取结果
	expectedPaths := []string{
		"./images/test1.png",
		"images/test2.jpg",
		"../images/test3.gif",
		"/absolute/path/image.png",
		"images\\test5.png",
		"images/test6.png?v=123",
		"./易立德产品数据管理软件(eRDCloud-PDM)部署手册_3.2.1.assets/1.png",
	}

	if len(paths) != len(expectedPaths) {
		t.Errorf("期望提取 %d 个路径，实际提取了 %d 个", len(expectedPaths), len(paths))
	}

	for i, expected := range expectedPaths {
		if i < len(paths) && paths[i] != expected {
			t.Errorf("路径 %d 不匹配: 期望 %s, 实际 %s", i+1, expected, paths[i])
		}
	}
}

func TestExtractImagesWithChinese(t *testing.T) {
	// 测试包含中文字符和特殊字符的图片路径
	content := `# 测试文档

![image-20250429112017900](./易立德产品数据管理软件(eRDCloud-PDM)部署手册_3.2.1.assets/image-20250429112017900.png)
![普通图片](./images/test.png)
![带参数图片](images/test.jpg?v=123)
![Windows路径](images\\test.gif)
![绝对路径](/absolute/path/image.png)
`

	// 提取图片路径
	paths := ExtractImages(content)

	// 验证提取结果
	expectedPaths := []string{
		"./易立德产品数据管理软件(eRDCloud-PDM)部署手册_3.2.1.assets/image-20250429112017900.png",
		"./images/test.png",
		"images/test.jpg?v=123",
		"images\\test.gif",
		"/absolute/path/image.png",
	}

	if len(paths) != len(expectedPaths) {
		t.Errorf("期望提取 %d 个路径，实际提取了 %d 个", len(expectedPaths), len(paths))
		t.Logf("实际提取的路径:")
		for i, path := range paths {
			t.Logf("  %d: %s", i+1, path)
		}
	}

	for i, expected := range expectedPaths {
		if i < len(paths) && paths[i] != expected {
			t.Errorf("路径 %d 不匹配: 期望 %s, 实际 %s", i+1, expected, paths[i])
		}
	}
}

func TestExtractImagesDebug(t *testing.T) {
	// 简单的调试测试
	content := `![image-20250429112017900](./易立德产品数据管理软件(eRDCloud-PDM)部署手册_3.2.1.assets/image-20250429112017900.png)`

	// 提取图片路径
	paths := ExtractImages(content)

	fmt.Printf("调试测试 - 提取到的路径数量: %d\n", len(paths))
	for i, path := range paths {
		fmt.Printf("调试测试 - 路径 %d: %s\n", i+1, path)
	}

	if len(paths) == 0 {
		t.Error("没有提取到任何图片路径")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
