@echo off

setlocal enabledelayedexpansion

::vars
set REST_EXEC_NAME=rest-api.exe
set GRAPH_EXEC_NAME=graph-api.exe

::param jump
if "%1"=="" goto all
if "%1"=="all" goto all
if "%1"=="setup" goto setup
if "%1"=="format" goto format
if "%1"=="docs" goto docs
if "%1"=="doc" goto docs
if "%1"=="checks" goto checks
if "%1"=="check" goto checks
if "%1"=="test" goto test
if "%1"=="tests" goto test
if "%1"=="test-graph" goto test-graph
if "%1"=="test-rest" goto test-rest
if "%1"=="test-shared" goto test-shared
if "%1"=="build" goto build
if "%1"=="build-rest" goto build-rest
if "%1"=="build-graph" goto build-graph
if "%1"=="generate" goto generate
if "%1"=="clean" goto clean

echo Unknown target: %1
echo Available targets: setup, format, docs, check, test, test-graph, test-rest, test-shared, generate, build, build-rest, build-graph, clean, all
exit /b 1

:all
call :setup
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

:format
echo Formatting...
go mod tidy
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
gofmt -w .
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
goimports -w .
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
echo Formatting done!
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
git diff --exit-code -- go.mod go.sum
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
go vet ./...
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
staticcheck ./...
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
set "GOFMT_FILES="
for /f "delims=" %%F in ('gofmt -l .') do (
    set "GOFMT_FILES=1"
)
if defined GOFMT_FILES exit /b 1
set "GOIMPORTS_FILES="
for /f "delims=" %%F in ('goimports -l .') do (
    set "GOIMPORTS_FILES=1"
)
if defined GOIMPORTS_FILES exit /b 1
echo Checks done!
echo.
exit /b 0

:test
echo Testing everything...
call :test-shared
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
call :test-rest
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
call :test-graph
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
echo Testing complete!
echo.
exit /b 0

:test-graph
go test ./graphql/... -count=1
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
echo GraphQL testing complete!
echo.
exit /b 0

:test-rest
echo Testing REST...
go test ./rest/... -count=1
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
echo REST testing complete!
echo.
exit /b 0

:test-shared
echo Testing shared...
go test ./shared/... -count=1
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
echo Shared testing complete!
echo.
exit /b 0

:build
call :docs
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
call :build-rest
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
call :build-graph
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
echo Build complete!
echo.
exit /b 0

:build-rest
echo Building...
go build -o %REST_EXEC_NAME% .\rest
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
echo REST build complete!
echo.
exit /b 0

:generate
echo Generating GraphQL execution layer...
pushd graphql
go get github.com/99designs/gqlgen@latest
if ERRORLEVEL 1 (
	popd
	exit /b %ERRORLEVEL%
)
go run github.com/99designs/gqlgen generate
if ERRORLEVEL 1 (
	popd
	exit /b %ERRORLEVEL%
)
popd
go mod tidy
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
echo GraphQL generation complete!
echo.
exit /b 0

:build-graph
echo Building GraphQL...
go build -o %GRAPH_EXEC_NAME% .\graphql
if ERRORLEVEL 1 exit /b %ERRORLEVEL%
echo GraphQL build complete!
echo.
exit /b 0

:clean
echo Cleaning...
if exist %REST_EXEC_NAME% del /f /q %REST_EXEC_NAME%
if exist rest\%REST_EXEC_NAME% del /f /q rest\%REST_EXEC_NAME%
if exist %GRAPH_EXEC_NAME% del /f /q %GRAPH_EXEC_NAME%
if exist graphql\%GRAPH_EXEC_NAME% del /f /q graphql\%GRAPH_EXEC_NAME%
echo Clean complete!
echo.
exit /b 0