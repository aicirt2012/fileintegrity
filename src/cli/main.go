package main

import (
	"os"

	"github.com/aicirt2012/fileintegrity/src/cli/cmd"
)

func main() {
	if err := cmd.Root().Execute(); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
