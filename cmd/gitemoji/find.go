package main

import (
	"io"
	"os"
	"strings"

	"github.com/coder/flog"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

func findCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "find",
		Long: "Find the best emoji for the commit message in stdin",
		Run: func(cmd *cobra.Command, _ []string) {
			midByte, err := fineTuneModelID.Read()
			if err != nil {
				if os.IsNotExist(err) {
					flog.Fatalf("fine-tuning model not found, run `gitemoji install`")
				}
				flog.Fatalf("read fine-tuning model id: %v", err)
			}

			stdin, err := io.ReadAll(os.Stdin)
			if err != nil {
				flog.Fatalf("read stdin: %v", err)
			}

			mid := strings.TrimSpace(string(midByte))

			client := newClient()
			resp, err := client.CreateCompletion(cmd.Context(), openai.CompletionRequest{
				Model: mid,
				// Model:     "davinci",
				MaxTokens: 2,
				Prompt:    string(stdin),
			})
			if err != nil {
				flog.Fatalf("create completion with %s: %v", mid, err)
			}
			flog.Infof("completion: %q", resp.Choices[0].Text)
		},
	}
}
