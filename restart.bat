@echo off
echo Stopping existing server...
taskkill /F /IM api.exe >nul 2>&1
taskkill /F /FI "WINDOWTITLE eq go run cmd/api/main.go" >nul 2>&1

echo Building and starting server...
go build -o api.exe cmd/api/main.go
start cmd /c "api.exe"

echo Server restarted!