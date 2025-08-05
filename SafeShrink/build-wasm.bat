@echo off
echo Building SafeShrink WASM module...

REM Check if Go is installed
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: Go is not installed or not in PATH
    echo Please install Go from https://golang.org/dl/
    pause
    exit /b 1
)

REM Set environment variables for WASM
set GOOS=js
set GOARCH=wasm

REM Create public directory if it doesn't exist
if not exist "public" mkdir public

REM Build the WASM file
echo Building WebAssembly module...
go build -o public/safeshrink.wasm main_wasm.go
if %errorlevel% neq 0 (
    echo Error: Failed to build WASM module
    pause
    exit /b 1
)

REM Find Go installation path and copy wasm_exec.js
for /f "tokens=*" %%i in ('go env GOROOT') do set GOROOT=%%i
copy "%GOROOT%\misc\wasm\wasm_exec.js" public\
if %errorlevel% neq 0 (
    echo Warning: Could not copy wasm_exec.js automatically
    echo Please manually copy wasm_exec.js from your Go installation
    echo Location: %GOROOT%\misc\wasm\wasm_exec.js
)

echo.
echo WASM module built successfully!
echo Files created:
echo   - public/safeshrink.wasm
echo   - public/wasm_exec.js
echo.
pause