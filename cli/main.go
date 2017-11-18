package main

import (
	"fmt"
	"os"

	"github.com/koki/control/cli/cmd/ctl"
)

func main() {
	if err := ctl.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
