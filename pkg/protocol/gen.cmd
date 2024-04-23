@echo off
setlocal

rem Define array elements
set "PROTO_NAMES=admin chat common"

rem Loop through each element in the array
for %%i in (%PROTO_NAMES%) do (
  protoc --go_out=plugins=grpc:./%%i --go_opt=module=github.com/openimsdk/chat/pkg/protocol/%%i %%i/%%i.proto
)

endlocal
