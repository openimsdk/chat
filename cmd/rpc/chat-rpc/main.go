package main

import (
	"github.com/openimsdk/chat/pkg/common/cmd"
	"github.com/openimsdk/tools/system/program"
)

func main() {
	if err := cmd.NewChatRpcCmd().Exec(); err != nil {
		program.ExitWithError(err)
	}
}
