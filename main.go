package main

import (
	"fmt"

	"github.com/bketelsen/bjk/cmd"
)

var (
	Version   = "N/A"
	BuildTime = "N/A"
)

func main() {
	fmt.Printf("Version: %s, Built: %s\n", Version, BuildTime)
	cmd.Execute()
}
