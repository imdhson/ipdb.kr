$ErrorActionPreference = "Stop"

$AppName = "ipdb.exe"
$ProcessName = "ipdb"

Write-Host "=== Building $AppName ==="
go build -o $AppName main.go
Write-Host "Build successful."

Write-Host "=== Deploying $AppName ==="

# Stop existing process if it's running
$existingProcess = Get-Process -Name $ProcessName -ErrorAction SilentlyContinue
if ($existingProcess) {
    Write-Host "Stopping existing instance (PID: $($existingProcess.Id))..."
    Stop-Process -Id $existingProcess.Id -Force
    Start-Sleep -Seconds 2
}

Write-Host "Starting new instance..."
# Run as a background job or separate process
Start-Process -FilePath ".\$AppName" -WindowStyle Hidden -RedirectStandardOutput "app.log" -RedirectStandardError "app.err.log"

Write-Host "Deployment complete! Application is running in the background."
Write-Host "Check app.log for output."
