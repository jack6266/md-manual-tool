package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// EnsureDir 确保目录存在，如果不存在则创建
func EnsureDir(path string) error {
	// 处理Windows路径
	path = filepath.Clean(path)

	// 处理UNC路径（网络路径）
	if strings.HasPrefix(path, `\\`) {
		// 对于网络路径，确保路径格式正确
		path = strings.ReplaceAll(path, "/", "\\")
	} else {
		// 对于本地路径，转换为绝对路径
		if !filepath.IsAbs(path) {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return fmt.Errorf("获取绝对路径失败: %v", err)
			}
			path = absPath
		}
	}

	// 检查路径是否有效
	// 只检查Windows不允许的字符，不检查中文等合法字符
	invalidChars := []string{"<", ">", "\"", "|", "?", "*"}
	for _, char := range invalidChars {
		if strings.Contains(path, char) {
			return fmt.Errorf("路径包含无效字符 '%s': %s", char, path)
		}
	}

	// 检查冒号（:）的使用是否合法
	// 只允许在盘符位置使用冒号
	if strings.Count(path, ":") > 1 || (strings.Contains(path, ":") && !strings.HasPrefix(path, `\\`) && !strings.HasPrefix(path, `//`)) {
		// 检查冒号是否在正确的位置（盘符后）
		parts := strings.Split(path, ":")
		if len(parts) > 1 && !strings.HasPrefix(parts[1], "\\") {
			return fmt.Errorf("冒号使用位置不正确: %s", path)
		}
	}

	// 如果是文件路径，则创建其父目录
	if strings.Contains(path, ".") {
		path = filepath.Dir(path)
	}

	// 处理长路径
	if len(path) > 260 {
		// 添加长路径前缀
		if !strings.HasPrefix(path, `\\?\`) {
			path = `\\?\` + path
		}
	}

	// 创建目录
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	return nil
}

// WriteFile 写入文件
func WriteFile(path string, content []byte) error {
	// 处理长路径
	if len(path) > 260 {
		// 添加长路径前缀
		if !strings.HasPrefix(path, `\\?\`) {
			path = `\\?\` + path
		}
	}

	// 确保目录存在
	if err := EnsureDir(path); err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(path, content, 0644)
}

// ReadFile 读取文件内容
func ReadFile(path string) ([]byte, error) {
	// 处理长路径
	if len(path) > 260 {
		// 添加长路径前缀
		if !strings.HasPrefix(path, `\\?\`) {
			path = `\\?\` + path
		}
	}
	return os.ReadFile(path)
}

// ExtractImages 从Markdown内容中提取图片路径
func ExtractImages(content string) []string {
	fmt.Println("ExtractImages收到的content内容如下:\n" + content)

	// 支持的图片格式
	imagePatterns := []string{
		`png`,  // PNG格式
		`jpg`,  // JPG格式
		`jpeg`, // JPEG格式
		`gif`,  // GIF格式
		`bmp`,  // BMP格式
		`webp`, // WebP格式
		`svg`,  // SVG格式
		`ico`,  // ICO格式
		`tiff`, // TIFF格式
		`tif`,  // TIF格式
	}

	var paths []string

	// 1. 匹配Markdown格式图片 ![alt](path)
	mdPattern := `(?s)!\[.*?\]\((.+?\.(?:` + strings.Join(imagePatterns, "|") + `)(?:\?[^)]*)?)\)`
	mdRe := regexp.MustCompile(mdPattern)
	mdMatches := mdRe.FindAllStringSubmatch(content, -1)

	for _, match := range mdMatches {
		if len(match) > 1 {
			// 保持原始路径格式，只去除首尾空格
			path := strings.TrimSpace(match[1])
			paths = append(paths, path)
		}
	}

	// 2. 匹配HTML格式图片 <img src="path" ... />
	// 更宽松的正则表达式，可以匹配各种格式的HTML图片标签
	htmlPattern := `<img\s+[^>]*?src=["']([^"']+?\.(?:` + strings.Join(imagePatterns, "|") + `)(?:[^"'>]*)?)["'][^>]*?>`
	htmlRe := regexp.MustCompile(htmlPattern)
	htmlMatches := htmlRe.FindAllStringSubmatch(content, -1)

	for _, match := range htmlMatches {
		if len(match) > 1 {
			// 保持原始路径格式，只去除首尾空格和多余的.png等后缀
			path := strings.TrimSpace(match[1])

			// 修复错误格式的路径，如果路径末尾有多余的.png等后缀
			for _, ext := range imagePatterns {
				doubleExt := "." + ext + "." + ext
				if strings.HasSuffix(path, doubleExt) {
					path = path[:len(path)-len("."+ext)]
				}
			}

			paths = append(paths, path)
		}
	}

	// 打印提取到的图片路径
	fmt.Printf("提取到的图片路径数量: %d\n", len(paths))
	for i, path := range paths {
		fmt.Printf("图片路径 %d: %s\n", i+1, path)
	}

	return paths
}

// normalizeImagePath 标准化图片路径
func normalizeImagePath(path string) string {
	// 移除URL参数（如 ?v=123）
	if idx := strings.Index(path, "?"); idx != -1 {
		path = path[:idx]
	}

	// 处理反斜杠，统一为斜杠
	path = strings.ReplaceAll(path, "\\", "/")

	// 移除开头的 ./
	if strings.HasPrefix(path, "./") {
		path = path[2:]
	}

	// 移除开头的 ../
	if strings.HasPrefix(path, "../") {
		path = path[3:]
	}

	return path
}

// resolveImagePath 解析图片路径（支持绝对路径和相对路径）
func resolveImagePath(imgPath, templatePath string) (string, error) {
	fmt.Printf("解析图片路径: %s (相对于模板: %s)\n", imgPath, templatePath)

	// 如果已经是绝对路径，直接返回
	if filepath.IsAbs(imgPath) {
		fmt.Printf("图片路径是绝对路径: %s\n", imgPath)
		return imgPath, nil
	}

	// 处理相对路径 - 保持原始格式，但尝试多种解析方式
	templateDir := filepath.Dir(templatePath)

	// 1. 直接拼接模板目录和图片路径
	resolvedPath := filepath.Join(templateDir, imgPath)
	if _, err := os.Stat(resolvedPath); err == nil {
		fmt.Printf("找到图片文件: %s\n", resolvedPath)
		return resolvedPath, nil
	}

	// 2. 处理以 ./ 开头的路径
	if strings.HasPrefix(imgPath, "./") {
		relativePath := imgPath[2:] // 移除 ./
		resolvedPath = filepath.Join(templateDir, relativePath)
		if _, err := os.Stat(resolvedPath); err == nil {
			fmt.Printf("找到图片文件 (./): %s\n", resolvedPath)
			return resolvedPath, nil
		}
	}

	// 3. 处理以 ../ 开头的路径
	if strings.HasPrefix(imgPath, "../") {
		relativePath := imgPath[3:] // 移除 ../
		parentDir := filepath.Dir(templateDir)
		resolvedPath = filepath.Join(parentDir, relativePath)
		if _, err := os.Stat(resolvedPath); err == nil {
			fmt.Printf("找到图片文件 (../): %s\n", resolvedPath)
			return resolvedPath, nil
		}
	}

	// 4. 尝试相对于当前工作目录
	currentDir, err := os.Getwd()
	if err == nil {
		resolvedPath = filepath.Join(currentDir, imgPath)
		if _, err := os.Stat(resolvedPath); err == nil {
			fmt.Printf("在当前目录找到图片文件: %s\n", resolvedPath)
			return resolvedPath, nil
		}
	}

	// 5. 尝试常见的图片目录
	commonDirs := []string{"images", "img", "assets", "pics", "pictures"}
	for _, dir := range commonDirs {
		// 相对于模板目录
		resolvedPath = filepath.Join(templateDir, dir, filepath.Base(imgPath))
		if _, err := os.Stat(resolvedPath); err == nil {
			fmt.Printf("在 %s 目录找到图片文件: %s\n", dir, resolvedPath)
			return resolvedPath, nil
		}

		// 相对于当前工作目录
		resolvedPath = filepath.Join(currentDir, dir, filepath.Base(imgPath))
		if _, err := os.Stat(resolvedPath); err == nil {
			fmt.Printf("在当前目录的 %s 目录找到图片文件: %s\n", dir, resolvedPath)
			return resolvedPath, nil
		}
	}

	// 6. 尝试处理Windows路径分隔符
	if strings.Contains(imgPath, "\\") {
		normalizedPath := strings.ReplaceAll(imgPath, "\\", "/")
		resolvedPath = filepath.Join(templateDir, normalizedPath)
		if _, err := os.Stat(resolvedPath); err == nil {
			fmt.Printf("找到图片文件 (Windows路径): %s\n", resolvedPath)
			return resolvedPath, nil
		}
	}

	return "", fmt.Errorf("无法找到图片文件: %s", imgPath)
}

// UpdateImagePaths 更新Markdown内容中的图片路径
func UpdateImagePaths(content string, mdPath string) string {
	// 获取Markdown文件名（不含扩展名）
	mdName := strings.TrimSuffix(filepath.Base(mdPath), filepath.Ext(mdPath))

	// 支持的图片格式
	imagePatterns := []string{
		`\.png`,  // PNG格式
		`\.jpg`,  // JPG格式
		`\.jpeg`, // JPEG格式
		`\.gif`,  // GIF格式
		`\.bmp`,  // BMP格式
		`\.webp`, // WebP格式
		`\.svg`,  // SVG格式
		`\.ico`,  // ICO格式
		`\.tiff`, // TIFF格式
		`\.tif`,  // TIF格式
	}

	updatedContent := content

	// 1. 更新Markdown格式的图片路径
	mdPattern := `!\[.*?\]\(([^\\)]+?(` + strings.Join(imagePatterns, "|") + `)(\\?[^\\)]*)?)\)`
	mdRe := regexp.MustCompile(mdPattern)

	updatedContent = mdRe.ReplaceAllStringFunc(updatedContent, func(match string) string {
		// 提取图片路径
		parts := mdRe.FindStringSubmatch(match)
		if len(parts) < 4 {
			return match
		}

		// 构建新的图片路径
		oldPath := parts[2] + parts[3]
		newPath := "./" + mdName + ".assets/" + filepath.Base(oldPath)

		// 返回更新后的图片标记
		return parts[1] + newPath + parts[4]
	})

	// 2. 更新HTML格式的图片路径
	// 更宽松的正则表达式，可以匹配各种格式的HTML图片标签
	htmlPattern := `(<img\s+[^>]*?src=)["']([^"']+?\.(?:` + strings.Join(imagePatterns, "|") + `)(?:[^"'>]*)?)(["'][^>]*?>)`
	htmlRe := regexp.MustCompile(htmlPattern)

	updatedContent = htmlRe.ReplaceAllStringFunc(updatedContent, func(match string) string {
		// 提取图片路径
		parts := htmlRe.FindStringSubmatch(match)
		if len(parts) < 4 {
			return match
		}

		// 处理路径中可能存在的错误格式
		path := parts[2]
		for _, ext := range imagePatterns {
			// 移除可能的双后缀
			doubleExt := ext + ext
			if strings.Contains(path, doubleExt) {
				path = strings.Replace(path, doubleExt, ext, 1)
			}
		}

		// 构建新的图片路径
		newPath := "./" + mdName + ".assets/" + filepath.Base(path)

		// 返回更新后的图片标记
		return parts[1] + "\"" + newPath + "\"" + parts[3]
	})

	return updatedContent
}

// CopyImagesFromTemplate 从模板文件复制图片到新目录并更新Markdown内容
func CopyImagesFromTemplate(templatePath, outputPath string, imagePaths []string, content string) (string, error) {
	fmt.Printf("开始处理图片复制...\n")
	fmt.Printf("模板文件路径: %s\n", templatePath)
	fmt.Printf("输出文件路径: %s\n", outputPath)
	fmt.Printf("图片路径数量: %d\n", len(imagePaths))

	// 获取输出文件名（不含扩展名）
	outputName := strings.TrimSuffix(filepath.Base(outputPath), filepath.Ext(outputPath))
	fmt.Printf("输出文件名: %s\n", outputName)

	// 创建图片目录
	imageDir := filepath.Join(filepath.Dir(outputPath), outputName+".assets")
	fmt.Printf("图片目录: %s\n", imageDir)
	err := EnsureDir(imageDir)
	if err != nil {
		return content, fmt.Errorf("创建图片目录失败: %v", err)
	}

	// 复制每个图片
	for i, imgPath := range imagePaths {
		fmt.Printf("\n处理图片 %d/%d: %s\n", i+1, len(imagePaths), imgPath)

		// 解析图片路径
		absImgPath, err := resolveImagePath(imgPath, templatePath)
		if err != nil {
			fmt.Printf("错误: 无法解析图片路径: %v\n", err)
			return content, fmt.Errorf("解析图片路径失败 %s: %v", imgPath, err)
		}

		// 处理长路径
		if len(absImgPath) > 260 {
			if !strings.HasPrefix(absImgPath, `\\?\`) {
				absImgPath = `\\?\` + absImgPath
			}
		}

		// 检查源文件是否存在
		if _, err := os.Stat(absImgPath); os.IsNotExist(err) {
			fmt.Printf("错误: 源图片文件不存在: %s\n", absImgPath)
			return content, fmt.Errorf("读取图片失败 %s: 文件不存在", imgPath)
		}

		// 读取源图片
		imgContent, err := os.ReadFile(absImgPath)
		if err != nil {
			fmt.Printf("错误: 读取图片文件失败: %v\n", err)
			return content, fmt.Errorf("读取图片失败 %s: %v", imgPath, err)
		}
		fmt.Printf("成功读取图片，大小: %d 字节\n", len(imgContent))

		// 写入新图片
		newImgPath := filepath.Join(imageDir, filepath.Base(imgPath))
		fmt.Printf("新图片路径: %s\n", newImgPath)
		err = WriteFile(newImgPath, imgContent)
		if err != nil {
			fmt.Printf("错误: 写入图片文件失败: %v\n", err)
			return content, fmt.Errorf("写入图片失败 %s: %v", newImgPath, err)
		}
		fmt.Printf("成功写入图片: %s\n", newImgPath)
	}

	// 更新Markdown内容中的图片路径
	fmt.Printf("\n更新Markdown内容中的图片路径...\n")
	updatedContent := UpdateImagePaths(content, outputPath)
	fmt.Printf("图片路径更新完成\n")

	return updatedContent, nil
}

// CopyImages 复制图片到新目录并更新Markdown内容
func CopyImages(mdPath string, imagePaths []string, content string) (string, error) {
	// 获取Markdown文件名（不含扩展名）
	mdName := strings.TrimSuffix(filepath.Base(mdPath), filepath.Ext(mdPath))

	// 创建图片目录
	imageDir := filepath.Join(filepath.Dir(mdPath), mdName+".assets")
	err := EnsureDir(imageDir)
	if err != nil {
		return content, fmt.Errorf("创建图片目录失败: %v", err)
	}

	// 复制每个图片
	for _, imgPath := range imagePaths {
		// 处理相对路径
		absImgPath := imgPath
		if !filepath.IsAbs(imgPath) {
			absImgPath = filepath.Join(filepath.Dir(mdPath), imgPath)
		}

		// 处理长路径
		if len(absImgPath) > 260 {
			if !strings.HasPrefix(absImgPath, `\\?\`) {
				absImgPath = `\\?\` + absImgPath
			}
		}

		// 读取源图片
		imgContent, err := os.ReadFile(absImgPath)
		if err != nil {
			return content, fmt.Errorf("读取图片失败 %s: %v", imgPath, err)
		}

		// 写入新图片
		newImgPath := filepath.Join(imageDir, filepath.Base(imgPath))
		err = WriteFile(newImgPath, imgContent)
		if err != nil {
			return content, fmt.Errorf("写入图片失败 %s: %v", newImgPath, err)
		}
	}

	// 更新Markdown内容中的图片路径
	updatedContent := UpdateImagePaths(content, mdPath)

	return updatedContent, nil
}
