@echo off
REM Usage: scripts\start.bat <PORT>
REM Example: scripts\start.bat 8080
if "%1"=="" (
    echo Usage: %0 ^<port^>
    exit /b 1
)
set PORT=%1
powershell -Command "(Get-Content nginx\nginx.conf.template) -replace '\${PORT}', '%PORT%' | Set-Content nginx\nginx.conf"
docker compose up --build --scale stock-service=3 -d
echo Stock Market running at http://localhost:%PORT%
