# Comprehensive Code Coverage Test Runner for All Backend Services
# This script runs tests and generates HTML coverage reports for each service

Write-Host "" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host "HMS Backend - Complete Coverage Report" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host "" -ForegroundColor Cyan

$services = @("auth-service", "building-service", "booking-service", "api-gateway")
$baseDir = "c:\Users\Dell\Desktop\HMS\Backend"

foreach ($service in $services) {
    Write-Host "" -ForegroundColor Yellow
    Write-Host "----------------------------------------" -ForegroundColor Yellow
    Write-Host "Testing: $service" -ForegroundColor Yellow
    Write-Host "----------------------------------------" -ForegroundColor Yellow
    
    $servicePath = Join-Path $baseDir $service
    
    if (Test-Path $servicePath) {
        Set-Location $servicePath
        
        Write-Host "Running tests..." -ForegroundColor Green
        
        # Run tests for all packages including main
        go test ./... -coverprofile=coverage.out -covermode=count
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "Tests passed" -ForegroundColor Green
            
            # Generate HTML coverage report
            go tool cover -html=coverage.out -o coverage.html
            Write-Host "Coverage report generated: coverage.html" -ForegroundColor Green
            
            # Display coverage summary
            Write-Host "" -ForegroundColor Cyan
            Write-Host "Coverage Summary:" -ForegroundColor Cyan
            go tool cover -func=coverage.out | Select-Object -Last 1
        } else {
            Write-Host "Tests failed for $service" -ForegroundColor Red
        }
    } else {
        Write-Host "Service directory not found: $servicePath" -ForegroundColor Red
    }
}

Write-Host "" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host "Opening all coverage reports..." -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host "" -ForegroundColor Cyan

# Open all HTML coverage reports in browser
foreach ($service in $services) {
    $coverageFile = Join-Path $baseDir "$service\coverage.html"
    if (Test-Path $coverageFile) {
        Start-Process $coverageFile
        Write-Host "Opened: $service/coverage.html" -ForegroundColor Green
    }
}

Write-Host "" -ForegroundColor Green
Write-Host "All coverage reports generated and opened!" -ForegroundColor Green
Write-Host "Check each service coverage.html file for detailed coverage." -ForegroundColor Cyan
