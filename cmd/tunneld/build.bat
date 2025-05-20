@REM set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-H windowsgui"

set GOOS=linux
set GOARCH=amd64
go build
