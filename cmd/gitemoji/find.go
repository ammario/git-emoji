package main

import (
	"io"
	"os"

	_ "embed"

	"github.com/coder/flog"
	"github.com/forPelevin/gomoji"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

//go:embed prelude.txt
var promptPrelude string

func findCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "find",
		Long: "Find the best emoji for the commit message in stdin",
		Run: func(cmd *cobra.Command, _ []string) {
			client := newClient()

			fineTunes, err := client.ListFineTunes(cmd.Context())
			if err != nil {
				flog.Fatalf("list fine-tunes: %v", err)
			}

			var (
				model      openai.FineTune
				foundModel bool
			)
		outer:
			for _, ft := range fineTunes.Data {
				if ft.FineTunedModel == "" {
					// This is either deleted or still training.
					continue
				}
				for _, file := range ft.TrainingFiles {
					if file.FileName != openAIFineTuneName {
						continue
					}
					model = ft
					foundModel = true
					break outer
				}
			}

			if !foundModel {
				flog.Fatalf("no model found, run `gitemoji install`")
			}

			stdin, err := io.ReadAll(os.Stdin)
			if err != nil {
				flog.Fatalf("read stdin: %v", err)
			}

			resp, err := client.CreateCompletion(cmd.Context(), openai.CompletionRequest{
				// Model:     model.FineTunedModel,
				Model:     "text-davinci-003",
				MaxTokens: 10,
				Prompt:    promptPrelude + "\nCommit: " + string(stdin),
			})
			if err != nil {
				flog.Fatalf("create completion with %+v: %v", model, err)
			}
			var best string
			for _, choice := range resp.Choices {
				flog.Infof("completion: %q", choice.Text)
				emojis := gomoji.FindAll(choice.Text)
				if len(emojis) > 0 {
					best = emojis[0].Character
					break
				}
			}
			if best == "" {
				flog.Errorf("no emoji found")
				os.Exit(20)
			}
		},
	}
}
