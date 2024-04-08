package main

import (
	"flag"
	"fmt"
	"github.com/openimsdk/chat/pkg/util"
	"github.com/openimsdk/chat/tools/mysql2mongo/internal"
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
