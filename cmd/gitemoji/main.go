package main

import (
	"github.com/coder/flog"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use: "gitemoji",
	}
	cmd.AddCommand(
		genFineTuningsCmd(),
		installCmd(),
		findCmd(),
	)

	err := cmd.Execute()
	if err != nil {
		flog.Fatalf("%v", err)
	}
}
