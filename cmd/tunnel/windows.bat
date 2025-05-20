@REM set ANDROID_SDK=H:\android_sdk
@REM set NDK_BIN=%ANDROID_SDK%\ndk\28.0.13004108\toolchains\llvm\prebuilt\windows-x86_64\bin

set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64
@REM set CC=%NDK_BIN%\aarch64-linux-android23-clang.cmd

@REM echo %CC%

go build -buildmode=c-shared -o H:/mydisk/deveplor/com.github/qq2383/tunnel_app/windows/runner/resources/tunnel.dll tunnel.go
