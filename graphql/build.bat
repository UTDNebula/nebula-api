@echo off

::vars
set EXEC_NAME=go-graph.exe

::param jump
if "%1"=="setup" goto setup
if "%1"=="checks" goto checks
if "%1"=="generate" goto generate
if "%1"=="build" goto build

:setup
echo Performing setup...
go install honnef.co/go/tools/cmd/staticcheck@latest && ^
go install golang.org/x/tools/cmd/goimports@latest
if ERRORLEVEL 1 exit /b %ERRORLEVEL% :: fail if error occurred
echo Setup done!
echo[
if "%1"=="setup" exit

:checks
echo Performing checks...
go mod tidy && ^
go vet ./... && ^
staticcheck ./... && ^
gofmt -w ./.. && ^
goimports -w ./..
if ERRORLEVEL 1 exit /b %ERRORLEVEL% :: fail if error occurred
echo Checks done!
echo[
if "%1"=="checks" exit

:generate
echo Generating GraphQL execution layer
go get github.com/99designs/gqlgen@latest && ^
go run github.com/99designs/gqlgen generate && ^
go mod tidy
if ERRORLEVEL 1 exit /b %ERRORLEVEL% :: fail if error occurred
echo Generation done!
echo[
if "%1"=="generate" exit

:test
echo Testing...
go test ./... -count=1
if ERRORLEVEL 1 exit /b %ERRORLEVEL% :: fail if error occurred
echo Testing complete!
echo[
if "%1"=="build" exit

:build
echo Building...
go build -o %EXEC_NAME% server.go
if ERRORLEVEL 1 exit /b %ERRORLEVEL% :: fail if error occurred
echo Build complete!
echo[
if "%1"=="build" exit