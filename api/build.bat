@echo off

::vars
set EXEC_NAME=go-api.exe

::setup
echo Performing setup...
go install honnef.co/go/tools/cmd/staticcheck@latest && ^
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/swaggo/swag/cmd/swag@latest
if ERRORLEVEL 1 exit /b %ERRORLEVEL% :: fail if error occurred
echo Setup done!
echo[

::docs
echo Generating docs...
swag init -g server.go

::checks
echo Performing checks...
go mod tidy && ^
go vet ./... && ^
staticcheck ./... && ^
gofmt -w ./.. && ^
goimports -w ./..
if ERRORLEVEL 1 exit /b %ERRORLEVEL% :: fail if error occurred
echo Checks done!
echo[

::build
echo Building...
go build -o %EXEC_NAME% server.go
if ERRORLEVEL 1 exit /b %ERRORLEVEL% :: fail if error occurred
echo Build complete!