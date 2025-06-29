package constants

// 默认配置
const (
	DefaultConfigFile = "config.yaml"
	DefaultOutputFile = "output.md"
)

// 用户提示消息
const (
	PromptTemplatePath = "请输入模板文件路径（如 templates/template.md）："
	PromptConfigPath   = "请输入配置文件路径（如 D:/config.yaml，直接回车则使用工具的当前目录的config.yaml）："
	PromptVersion      = "请输入版本号（如 1.0.1）："
)

// 成功消息
const (
	MsgDefaultConfigUsed = "使用默认配置文件：%s"
	MsgVersionDetected   = "检测到模板文件中的版本号：%s"
	MsgVersionAdded      = "版本参数已添加：%s"
	MsgFileGenerated     = "文件生成成功！输出路径：%s"
)

// 错误消息
const (
	ErrReadTemplatePath = "读取模板文件路径失败: %v"
	ErrReadConfigPath   = "读取配置文件路径失败: %v"
	ErrReadVersion      = "读取版本号失败: %v"
	ErrFileNotExist     = "%s不存在: %s\n当前工作目录: %s"
	ErrReadConfig       = "读取配置文件失败: %v"
	ErrProcessFailed    = "处理失败: %v"
	ErrCollectInputs    = "收集输入失败: %v"
	ErrValidateInputs   = "验证输入失败: %v"
	ErrLoadConfig       = "加载配置失败: %v"
	ErrProcessDocument  = "处理文档失败: %v"
)

// 文件类型
const (
	FileTypeTemplate = "模板文件"
	FileTypeConfig   = "配置文件"
)

// 版本号正则表达式
const (
	VersionRegexPattern = `_(\d+\.\d+\.\d+)\.md$`
)
