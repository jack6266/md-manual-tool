package main

import (
	"bufio"
	"fmt"
	"md-manual-tool/pkg/config"
	"md-manual-tool/pkg/constants"
	"md-manual-tool/pkg/document"
	"md-manual-tool/pkg/input"
	"md-manual-tool/pkg/ui"
	"md-manual-tool/pkg/validator"
	"os"
)

// Application 应用程序结构体
type Application struct {
	collector    *input.Collector
	validator    *validator.Validator
	configMgr    *config.Manager
	docProcessor *document.Processor
	ui           *ui.Interface
}

// NewApplication 创建新的应用程序实例
func NewApplication() *Application {
	reader := bufio.NewReader(os.Stdin)

	return &Application{
		collector:    input.NewCollector(reader),
		validator:    validator.NewValidator(),
		configMgr:    config.NewManager(),
		docProcessor: document.NewProcessor(),
		ui:           ui.NewInterface(),
	}
}

// Run 运行应用程序
func (app *Application) Run() error {
	// 1. 收集用户输入
	inputData, err := app.collectInputs()
	if err != nil {
		return fmt.Errorf(constants.ErrCollectInputs, err)
	}

	// 2. 验证输入
	if err := app.validateInputs(inputData); err != nil {
		return fmt.Errorf(constants.ErrValidateInputs, err)
	}

	// 3. 加载和处理配置
	configData, err := app.loadConfig(inputData)
	if err != nil {
		return fmt.Errorf(constants.ErrLoadConfig, err)
	}

	// 4. 处理文档
	if err := app.processDocument(configData); err != nil {
		return fmt.Errorf(constants.ErrProcessDocument, err)
	}

	// 5. 显示成功信息
	app.ui.ShowSuccess(configData.OutputPath)
	return nil
}

// collectInputs 收集用户输入
func (app *Application) collectInputs() (*input.InputData, error) {
	return app.collector.CollectAll()
}

// validateInputs 验证输入
func (app *Application) validateInputs(inputData *input.InputData) error {
	// 验证文件存在性
	result := app.validator.ValidateInputs(inputData.TemplatePath, inputData.ConfigPath)
	if !result.IsValid {
		app.ui.ShowValidationErrors(result.Errors)
		return fmt.Errorf("输入验证失败")
	}

	// 验证版本号格式
	if err := app.validator.ValidateVersionFormat(inputData.Version); err != nil {
		return err
	}

	return nil
}

// loadConfig 加载配置
func (app *Application) loadConfig(inputData *input.InputData) (*config.ConfigData, error) {
	return app.configMgr.LoadAndProcessConfig(
		inputData.ConfigPath,
		inputData.TemplatePath,
		inputData.Version,
	)
}

// processDocument 处理文档
func (app *Application) processDocument(configData *config.ConfigData) error {
	return app.docProcessor.ProcessDocument(configData)
}

func main() {
	app := NewApplication()
	if err := app.Run(); err != nil {
		fmt.Printf("错误：%v\n", err)
		os.Exit(1)
	}
}
