# SafeShrink WASM Build Script

Write-Host "Building WASM module..." -ForegroundColor Green

# Check Go
try {
    go version | Out-Null
    Write-Host "Go found" -ForegroundColor Cyan
} catch {
    Write-Host "Go not found. Install from https://golang.org/dl/" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

# Set WASM environment
$env:GOOS = "js"
$env:GOARCH = "wasm"

# Create public folder
if (!(Test-Path "public")) {
    New-Item -ItemType Directory -Path "public" | Out-Null
}

# Build WASM
try {
    go build -o public/safeshrink.wasm main_wasm.go
    Write-Host "WASM built successfully" -ForegroundColor Green
} catch {
    Write-Host "Build failed" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

# Copy wasm_exec.js
try {
    $goRoot = go env GOROOT
    $wasmExecPath = Join-Path $goRoot "misc\wasm\wasm_exec.js"
    
    if (Test-Path $wasmExecPath) {
        Copy-Item $wasmExecPath "public\"
        Write-Host "Files copied" -ForegroundColor Green
    } else {
        Write-Host "Copy wasm_exec.js manually from Go installation" -ForegroundColor Yellow
    }
} catch {
    Write-Host "Copy wasm_exec.js manually from Go installation" -ForegroundColor Yellow
}

Write-Host "`nBuild complete! Run: npm run dev" -ForegroundColor Green
Read-Host "Press Enter to exit"