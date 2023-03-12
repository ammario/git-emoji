package main

import (
	"io"
	"os"

	"github.com/coder/flog"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

const promptPrelude = "Print the emoji that best describes the following commit message. Print nothing but the emoji.\n"

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
				// Model: model.FineTunedModel,
				Model:       "text-davinci-003",
				MaxTokens:   3,
				Temperature: 0.9,
				Prompt:      promptPrelude + string(stdin),
			})
			if err != nil {
				flog.Fatalf("create completion with %+v: %v", model, err)
			}
			flog.Infof("completion: %q", resp.Choices[0].Text)
		},
	}
}
