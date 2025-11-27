# PowerShell script to generate Go code from proto files

Write-Host "Generating gRPC code from proto files..." -ForegroundColor Green

# Auth Service
Write-Host "`nGenerating auth.proto..." -ForegroundColor Yellow
protoc --go_out=..\auth-service --go_opt=paths=source_relative `
    --go-grpc_out=..\auth-service --go-grpc_opt=paths=source_relative `
    auth.proto

if ($LASTEXITCODE -ne 0) {
    Write-Host "Error generating auth.proto" -ForegroundColor Red
    exit 1
}

# Building Service
Write-Host "Generating building.proto..." -ForegroundColor Yellow
protoc --go_out=..\building-service --go_opt=paths=source_relative `
    --go-grpc_out=..\building-service --go-grpc_opt=paths=source_relative `
    building.proto

if ($LASTEXITCODE -ne 0) {
    Write-Host "Error generating building.proto" -ForegroundColor Red
    exit 1
}

Write-Host "`nProto files generated successfully!" -ForegroundColor Green
