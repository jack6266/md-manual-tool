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

// ExtractImages 从Markdown内容中提取图片路径
func ExtractImages(content string) []string {
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

	// 构建正则表达式
	pattern := `!\[.*?\]\((.*?(` + strings.Join(imagePatterns, "|") + `))\)`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(content, -1)

	var paths []string
	for _, match := range matches {
		if len(match) > 1 {
			// 处理相对路径中的反斜杠
			path := strings.ReplaceAll(match[1], "\\", "/")
			paths = append(paths, path)
		}
	}
	return paths
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

	// 构建正则表达式
	pattern := `(!\[.*?\]\()(.*?)(` + strings.Join(imagePatterns, "|") + `)(\))`
	re := regexp.MustCompile(pattern)

	// 替换图片路径
	updatedContent := re.ReplaceAllStringFunc(content, func(match string) string {
		// 提取图片路径
		parts := re.FindStringSubmatch(match)
		if len(parts) < 4 {
			return match
		}

		// 构建新的图片路径
		oldPath := parts[2] + parts[3]
		newPath := "./" + mdName + ".assets/" + filepath.Base(oldPath)

		// 返回更新后的图片标记
		return parts[1] + newPath + parts[4]
	})

	return updatedContent
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
