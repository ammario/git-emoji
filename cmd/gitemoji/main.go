package main

import (
	"os"

	"github.com/coder/flog"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

func newClient() *openai.Client {
	const keyEnv = "OPENAI_API_KEY"
	key := os.Getenv(keyEnv)
	if key == "" {
		flog.Fatalf("$%v is not set", keyEnv)
	}
	return openai.NewClient(key)
}

func main() {
	cmd := &cobra.Command{
		Use: "gitemoji",
	}
	cmd.AddCommand(
		amendCmd(),
	)

	err := cmd.Execute()
	if err != nil {
		flog.Fatalf("%v", err)
	}
}
