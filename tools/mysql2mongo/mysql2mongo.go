package main

import (
	"flag"
	"fmt"

	"github.com/OpenIMSDK/chat/pkg/util"
	"github.com/OpenIMSDK/chat/tools/mysql2mongo/internal"
)

func main() {
	var path string
	flag.StringVar(&path, "c", "", "path config file")
	flag.Parse()
	if err := internal.Main(path); err != nil {
		util.ExitWithError(err)
	}
	fmt.Println("chat mysql2mongo success!")
	return
}
