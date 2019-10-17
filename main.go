package main

import (
	"os"

	"github.com/defus/springboothystrixbeat/cmd"

	_ "github.com/defus/springboothystrixbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
