@echo off
chcp 65001 > nul

echo 开始构建 md-manual-tool...

:: 清理旧的构建文件
if exist md-manual-tool.exe (
    del md-manual-tool.exe
    echo 已删除旧的构建文件
)

:: 构建程序
echo 正在构建程序...
go build -o md-manual-tool.exe

if %ERRORLEVEL% NEQ 0 (
    echo 构建失败！
    exit /b 1
)

echo 构建成功！
echo 可执行文件位置: %CD%\md-manual-tool.exe 