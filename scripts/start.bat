@echo off
REM Usage: scripts\start.bat <PORT>
REM Example: scripts\start.bat 8080
if "%~1"=="" (
    echo Usage: %~0 ^<port^>
    exit /b 1
)
set PORT=%~1

REM Generate nginx.conf from template (replaces ${PORT})
powershell -NoProfile -Command ^
  "(Get-Content nginx\nginx.conf.template -Raw) -replace '\$\{PORT\}', '%PORT%' | Set-Content nginx\nginx.conf -NoNewline"

docker compose up --build -d

echo.
echo [OK] Stock Market running at http://localhost:%PORT%
echo      Instances: stock-service-1, stock-service-2, stock-service-3
echo      Shared state: Redis
