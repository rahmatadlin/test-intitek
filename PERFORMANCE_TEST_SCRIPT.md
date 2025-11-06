# 🧪 Script untuk Testing Performance

Script-script untuk mengukur performance Tauri vs Wails.

## 📊 Memory & CPU Monitoring Script

### Windows PowerShell Script

Simpan sebagai `measure-performance.ps1`:

```powershell
# measure-performance.ps1
# Usage: .\measure-performance.ps1 "app-name.exe"

param(
    [string]$ProcessName = "warehouse-management-frontend.exe"
)

Write-Host "Monitoring $ProcessName..." -ForegroundColor Green
Write-Host "Press Ctrl+C to stop" -ForegroundColor Yellow
Write-Host ""

$headers = "Time", "CPU%", "Memory(MB)", "Threads"
$results = @()

while ($true) {
    $process = Get-Process -Name $ProcessName -ErrorAction SilentlyContinue
    
    if ($process) {
        $cpu = $process.CPU
        $memory = [math]::Round($process.WorkingSet64 / 1MB, 2)
        $threads = $process.Threads.Count
        $time = Get-Date -Format "HH:mm:ss"
        
        $results += [PSCustomObject]@{
            Time = $time
            CPU = $cpu
            Memory = $memory
            Threads = $threads
        }
        
        Write-Host "$time - CPU: $cpu% | Memory: $memory MB | Threads: $threads"
        
        Start-Sleep -Seconds 2
    } else {
        Write-Host "Process not found. Waiting..." -ForegroundColor Red
        Start-Sleep -Seconds 2
    }
}

# Export results
$results | Export-Csv -Path "performance-results.csv" -NoTypeInformation
Write-Host "Results saved to performance-results.csv" -ForegroundColor Green
```

### Usage:
```powershell
# Untuk Tauri
.\measure-performance.ps1 "warehouse-management-frontend.exe"

# Untuk Wails
.\measure-performance.ps1 "warehouse-management.exe"
```

---

## ⏱️ Build Time Measurement Script

### Windows PowerShell

```powershell
# measure-build-time.ps1
# Usage: .\measure-build-time.ps1 "tauri" atau "wails"

param(
    [Parameter(Mandatory=$true)]
    [ValidateSet("tauri","wails")]
    [string]$Framework
)

$stopwatch = [System.Diagnostics.Stopwatch]::StartNew()

Write-Host "Measuring $Framework build time..." -ForegroundColor Green
Write-Host "Starting build..." -ForegroundColor Yellow

if ($Framework -eq "tauri") {
    cd frontend
    npm run tauri:build
} else {
    wails build
}

$stopwatch.Stop()

$buildTime = $stopwatch.Elapsed

Write-Host ""
Write-Host "Build completed!" -ForegroundColor Green
Write-Host "Total time: $($buildTime.TotalSeconds) seconds" -ForegroundColor Cyan
Write-Host "Total time: $($buildTime.TotalMinutes) minutes" -ForegroundColor Cyan

# Save to file
$result = @{
    Framework = $Framework
    TotalSeconds = $buildTime.TotalSeconds
    TotalMinutes = $buildTime.TotalMinutes
    Timestamp = Get-Date
} | ConvertTo-Json

Add-Content -Path "build-time-results.json" -Value $result
```

---

## 📦 File Size Measurement Script

### Windows PowerShell

```powershell
# measure-file-size.ps1
# Usage: .\measure-file-size.ps1 "path-to-exe"

param(
    [Parameter(Mandatory=$true)]
    [string]$FilePath
)

if (Test-Path $FilePath) {
    $file = Get-Item $FilePath
    $sizeBytes = $file.Length
    $sizeMB = [math]::Round($sizeBytes / 1MB, 2)
    $sizeKB = [math]::Round($sizeBytes / 1KB, 2)
    
    Write-Host "File: $($file.Name)" -ForegroundColor Green
    Write-Host "Path: $($file.FullName)" -ForegroundColor Yellow
    Write-Host "Size: $sizeMB MB ($sizeKB KB)" -ForegroundColor Cyan
    
    # Get folder size
    $folder = $file.Directory
    $folderSize = (Get-ChildItem -Path $folder.FullName -Recurse | 
        Measure-Object -Property Length -Sum).Sum
    $folderSizeMB = [math]::Round($folderSize / 1MB, 2)
    
    Write-Host "Folder Size: $folderSizeMB MB" -ForegroundColor Cyan
    
    # Save to JSON
    $result = @{
        FileName = $file.Name
        FilePath = $file.FullName
        SizeBytes = $sizeBytes
        SizeMB = $sizeMB
        SizeKB = $sizeKB
        FolderSizeMB = $folderSizeMB
        Timestamp = Get-Date
    } | ConvertTo-Json
    
    Add-Content -Path "file-size-results.json" -Value $result
} else {
    Write-Host "File not found: $FilePath" -ForegroundColor Red
}
```

---

## 🚀 Startup Time Measurement Script

### Windows PowerShell

```powershell
# measure-startup-time.ps1
# Usage: .\measure-startup-time.ps1 "path-to-exe"

param(
    [Parameter(Mandatory=$true)]
    [string]$ExePath
)

Write-Host "Measuring startup time for: $ExePath" -ForegroundColor Green
Write-Host "Starting application..." -ForegroundColor Yellow

$stopwatch = [System.Diagnostics.Stopwatch]::StartNew()

# Start process
$process = Start-Process -FilePath $ExePath -PassThru

# Wait until window is visible
$windowVisible = $false
$maxWaitTime = 30 # seconds
$elapsed = 0

while (-not $windowVisible -and $elapsed -lt $maxWaitTime) {
    Start-Sleep -Milliseconds 100
    $elapsed += 0.1
    
    try {
        $proc = Get-Process -Id $process.Id -ErrorAction SilentlyContinue
        if ($proc -and $proc.MainWindowHandle -ne 0) {
            $windowVisible = $true
        }
    } catch {
        # Process might not have window handle yet
    }
}

$stopwatch.Stop()

if ($windowVisible) {
    Write-Host "Window appeared!" -ForegroundColor Green
    Write-Host "Startup time: $($stopwatch.ElapsedMilliseconds) ms" -ForegroundColor Cyan
    Write-Host "Startup time: $($stopwatch.Elapsed.TotalSeconds) seconds" -ForegroundColor Cyan
    
    # Save result
    $result = @{
        ExePath = $ExePath
        StartupTimeMs = $stopwatch.ElapsedMilliseconds
        StartupTimeSeconds = $stopwatch.Elapsed.TotalSeconds
        Timestamp = Get-Date
    } | ConvertTo-Json
    
    Add-Content -Path "startup-time-results.json" -Value $result
    
    # Wait a bit, then close
    Write-Host "Closing application in 5 seconds..." -ForegroundColor Yellow
    Start-Sleep -Seconds 5
    Stop-Process -Id $process.Id -Force
} else {
    Write-Host "Timeout waiting for window to appear" -ForegroundColor Red
    Stop-Process -Id $process.Id -Force
}
```

---

## 📋 Quick Test Checklist

### Setup Test
```powershell
# 1. Clean start
Remove-Item -Recurse -Force node_modules, dist, target -ErrorAction SilentlyContinue

# 2. Measure setup time
Measure-Command {
    npm install
    # ... setup commands
}
```

### Build Test
```powershell
# 1. Clean build
Remove-Item -Recurse -Force dist, target -ErrorAction SilentlyContinue

# 2. Measure build time
.\measure-build-time.ps1 "tauri"
```

### Runtime Test
```powershell
# 1. Start app
Start-Process "path\to\app.exe"

# 2. Start monitoring
.\measure-performance.ps1 "app-name.exe"

# 3. Perform operations
# - Load data
# - Navigate pages
# - Heavy operations

# 4. Stop monitoring (Ctrl+C)
# 5. Close app
```

---

## 📸 Screenshot Checklist

### Terminal Outputs
- [ ] Setup command output
- [ ] Build command output
- [ ] Error messages
- [ ] Performance measurement output

### File System
- [ ] Folder structure
- [ ] File properties (size)
- [ ] Bundle contents

### Runtime
- [ ] Task Manager (Memory/CPU)
- [ ] App window (idle state)
- [ ] App window (loaded state)
- [ ] DevTools Performance tab

### Comparison
- [ ] Side-by-side screenshots
- [ ] Comparison tables
- [ ] Charts/graphs

---

**Tips:**
- Gunakan screenshot tool dengan timestamp
- Simpan screenshot dengan nama deskriptif: `tauri-build-time.png`, `wails-memory-usage.png`
- Buat folder terpisah untuk setiap framework: `screenshots/tauri/`, `screenshots/wails/`
- Gunakan video recording untuk hot reload demonstration

