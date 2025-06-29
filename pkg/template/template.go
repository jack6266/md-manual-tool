package template

import (
	"bytes"
	"io/ioutil"
	"text/template"
)

// Render 渲染模板
func Render(templatePath string, variables map[string]string) ([]byte, error) {
	// 读取模板文件
	templateContent, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return nil, err
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

	return result.Bytes(), nil
}
