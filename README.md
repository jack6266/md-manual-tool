# Markdown模板渲染工具

这是一个使用Go语言实现的Markdown模板渲染工具，可以根据模板文件和配置文件生成最终的Markdown文档。

## 环境要求

- Go 1.21 或更高版本
- Windows 操作系统

## 构建步骤

### 方法一：使用构建脚本（推荐）

1. 打开命令提示符或PowerShell
2. 进入项目目录
3. 运行构建脚本：
   ```bash
   build.bat
   ```

### 方法二：手动构建

1. 打开命令提示符或PowerShell
2. 进入项目目录
3. 运行构建命令：
   ```bash
   go build -o md-manual-tool.exe
   ```

## 使用方法

### 基本用法
1. 在 `templates` 目录下创建模板文件，使用 `{{.变量名}}` 的形式定义变量
2. 在 `configs` 目录下创建配置文件，定义变量的值
3. 运行程序：
   ```bash
   md-manual-tool.exe
   ```
4. 生成的Markdown文件将保存在 `output` 目录下

### 命令行参数
程序支持以下命令行参数：

- `-template`：指定模板文件路径（默认为 `templates/template.md`）
- `-config`：指定配置文件路径（默认为 `configs/config.yaml`）
- `-output`：指定输出文件路径（默认为 `output/result.md`）

使用示例：
```bash
# 使用默认路径
md-manual-tool.exe

# 指定所有参数
md-manual-tool.exe -template my-template.md -config my-config.yaml -output my-output.md

# 查看帮助信息
md-manual-tool.exe -help
```

## 项目结构

```
.
├── build.bat           # 构建脚本
├── md-manual-tool.exe  # 可执行文件
├── go.mod              # Go模块文件
├── main.go             # 主程序文件
├── pkg/                # 核心包
│   ├── config/         # 配置模块
│   │   └── config.go
│   ├── template/       # 模板模块
│   │   └── template.go
│   ├── utils/          # 工具模块
│   │   └── utils.go
│   └── processor/      # 处理器模块
│       └── processor.go
├── templates/          # 模板文件目录
│   └── template.md
├── configs/            # 配置文件目录
│   └── config.yaml
└── output/             # 输出文件目录
    └── result.md
```

## 支持的图片格式

程序支持以下图片格式：
- PNG (`.png`)
- JPG (`.jpg`, `.jpeg`)
- GIF (`.gif`)
- BMP (`.bmp`)
- WebP (`.webp`)
- SVG (`.svg`)
- ICO (`.ico`)
- TIFF (`.tiff`, `.tif`)

## 示例

模板文件示例：
```markdown
# {{.title}}

## 项目简介
{{.description}}

## 图片示例
![示例图片](./images/example.png)
```

配置文件示例：
```yaml
title: 示例项目
description: 这是一个示例描述
```

## 代码架构

### 模块化设计
项目采用模块化设计，主要包含以下模块：

1. **config模块**：负责配置文件的读取和解析
2. **template模块**：负责模板的渲染
3. **utils模块**：提供文件操作和图片处理工具函数
4. **processor模块**：协调各个模块，处理整个工作流程

### 主要特性
- 支持Windows长路径（超过260字符）
- 自动处理图片复制和路径更新
- 支持多种图片格式
- 模块化设计，易于维护和扩展
- 完善的错误处理机制 