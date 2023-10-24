@echo off

::vars
set EXEC_NAME=go-api.exe

::setup
echo Performing setup...
go install honnef.co/go/tools/cmd/staticcheck@latest
go install golang.org/x/tools/cmd/goimports@latest
echo Setup done!
echo[

::checks
echo Performing checks...
go mod tidy
go vet ./... 
staticcheck ./...
gofmt -w ./..
goimports -w ./..
echo Checks done!
echo[

::build
echo Building...
go build -o %EXEC_NAME% server.go
echo Build complete!