package main

import (
	"flag"
	"github.com/OpenIMSDK/chat/tools/mysql2mongo/internal"
	"log"
)

func main() {
	var path string
	flag.StringVar(&path, "c", "", "path config file")
	flag.Parse()
	log.SetFlags(log.Llongfile | log.Ldate | log.Ltime)
	if err := internal.Main(path); err != nil {
		log.Fatal("chat mysql2mongo error", err)
		return
	}
	log.Println("chat mysql2mongo success!")
	return
}
