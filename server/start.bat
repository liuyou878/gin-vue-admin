@echo off
cd /d D:\afg\Alpha_ERP\server
go build -o server.exe main.go
if %errorlevel%==0 (
    echo 编译成功，启动程序...
    start server.exe
) else (
    echo 编译失败，请检查代码
    pause
)