package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	_ "embed"

	"github.com/coder/flog"
	"github.com/fatih/color"
	"github.com/forPelevin/gomoji"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

//go:embed prelude.txt
var promptPrelude string

func readCommitFromCwd() (string, error) {
	var buf bytes.Buffer
	c := exec.Command("git", "log", "--pretty=format:%s", "-1")
	c.Stdout = &buf
	err := c.Run()
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func amendName(s string) error {
	c := exec.Command("git", "commit", "--amend", "-m", s)
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}

func amendCmd() *cobra.Command {
	var (
		readStdin bool
		dryRun    bool
	)
	cmd := &cobra.Command{
		Use:  "amend",
		Long: "Find the best emoji for the commit message in stdin",
		Run: func(cmd *cobra.Command, _ []string) {
			client := newClient()

			var commitMessage string
			if readStdin {
				byt, err := io.ReadAll(os.Stdin)
				if err != nil {
					flog.Fatalf("read stdin: %v", err)
				}
				commitMessage = string(byt)
			} else {
				var err error
				commitMessage, err = readCommitFromCwd()
				if err != nil {
					flog.Fatalf("read commit: %v", err)
				}
			}

			// Remove existing emojis to make it easy to retry.
			commitMessage = gomoji.RemoveEmojis(commitMessage)

			commitMessage = strings.TrimSpace(commitMessage)

			resp, err := client.CreateCompletion(cmd.Context(), openai.CompletionRequest{
				Model:     "text-davinci-003",
				MaxTokens: 10,
				Prompt:    promptPrelude + "\nCommit: " + commitMessage,
			})
			if err != nil {
				flog.Fatalf("create completion: %v", err)
			}
			var best string
			for _, choice := range resp.Choices {
				// flog.Infof("completion: %q", choice.Text)
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

			newCommitMessage := best + " " + commitMessage

			if dryRun {
				color.Blue("Dry run, would have amended to: \n")
			} else {
				color.Magenta("> git commit --amend \n")
				err := amendName(newCommitMessage)
				if err != nil {
					flog.Fatalf("amend: %v", err)
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", newCommitMessage)
		},
	}

	cmd.Flags().BoolVar(&readStdin, "read-stdin", false, "read commit message from stdin")
	cmd.Flags().BoolVar(&dryRun, "dry", false, "don't actually amend the commit, just print the new message")

	return cmd
}
