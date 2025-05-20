set ANDROID_SDK=H:\android_sdk
set NDK_BIN=%ANDROID_SDK%\ndk\23.1.7779620\toolchains\llvm\prebuilt\windows-x86_64\bin
set ANDROID_OUT=H:\mydisk\deveplor\com.github\qq2383\tunnel_app\android\app\src\main\jniLibs

set CGO_ENABLED=1

set GOOS=android
set GOARCH=arm64
set CC=%NDK_BIN%\aarch64-linux-android26-clang.cmd
go build -buildmode=c-shared -o %ANDROID_OUT%\arm64-v8a\tunnel.so tunnel.go
@REM go build -buildmode=c-shared -o %ANDROID_OUT%\tunnel.so tunnel.go

set GOOS=android
set GOARM=7
set GOARCH=arm
set CC=%NDK_BIN%\armv7a-linux-androideabi26-clang.cmd
go build -buildmode=c-shared -o %ANDROID_OUT%\armeabi-v7a\tunnel.so tunnel.go

set GOOS=android
set GOARCH=386
set CC=%NDK_BIN%\i686-linux-android26-clang.cmd
go build -buildmode=c-shared -o %ANDROID_OUT%\x86\tunnel.so tunnel.go

set GOOS=android
set GOARCH=amd64
set CC=%NDK_BIN%\x86_64-linux-android26-clang.cmd
go build -buildmode=c-shared -o %ANDROID_OUT%\x86_64\tunnel.so tunnel.go

@REM gomobile bind -o .\android\tunnel.aar -target=android -androidapi=21 .\android