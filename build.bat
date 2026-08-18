@echo off
setlocal enabledelayedexpansion

::vars
set EXEC_NAME=go-api.exe

::param jump
if "%1"=="" goto all
if "%1"=="all" goto all
if "%1"=="setup" goto setup
if "%1"=="docs" goto docs
if "%1"=="doc" goto docs
if "%1"=="checks" goto checks
if "%1"=="check" goto checks
if "%1"=="test" goto test
if "%1"=="tests" goto test
if "%1"=="build" goto build
if "%1"=="clean" goto clean

echo Unknown target: %1
echo Available targets: setup, docs, check, test, build, clean, all
exit /b 1

:all
call :setup
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
call :docs
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
call :checks
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
call :test
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
call :build
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
exit /b 0

:setup
echo Performing setup...
go install honnef.co/go/tools/cmd/staticcheck@latest
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
go install golang.org/x/tools/cmd/goimports@latest
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
go install github.com/swaggo/swag/cmd/swag@latest
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
echo Setup done!
echo.
exit /b 0

:docs
echo Generating docs...
swag fmt -d rest
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
swag init -d rest -g server.go -o rest\docs --outputTypes yaml,go
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
echo Docs generated!
echo.
exit /b 0

:checks
echo Performing checks...
go mod tidy
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
go vet ./...
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
staticcheck ./...
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
gofmt -w .
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
goimports -w .
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
echo Checks done!
echo.
exit /b 0

:test
echo Testing...
go test ./... -count=1
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
echo Testing complete!
echo.
exit /b 0

:build
echo Generating docs...
swag fmt -d rest
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
swag init -d rest -g server.go -o rest\docs --outputTypes yaml,go
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
echo Building...
go build -o %EXEC_NAME% .\rest
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
echo Build complete!
echo.
exit /b 0

:clean
echo Cleaning...
if exist %EXEC_NAME% del /f /q %EXEC_NAME%
if exist rest\%EXEC_NAME% del /f /q rest\%EXEC_NAME%
echo Clean complete!
echo.
exit /b 0