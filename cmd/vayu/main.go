package main

import (
	"fmt"
	"os"

	"github.com/Rohanraj123/vayu/cmd/vayu/commands"
	"github.com/spf13/cobra"
)

func main() {
	cmd, err := setup()
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func setup() (*cobra.Command, error) {
	cmd := commands.RootCommand()

	return cmd, nil
}
