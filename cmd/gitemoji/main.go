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
	"github.com/manifoldco/promptui"
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

func main() {
	var (
		readStdin     bool
		dryRun        bool
		adventureMode bool
	)
	cmd := &cobra.Command{
		Use:   "gitemoji",
		Short: "Add an emoji to your last commit message",
		Run: func(cmd *cobra.Command, args []string) {
			client := newClient()

			var originalCommitMessage string
			if readStdin {
				byt, err := io.ReadAll(os.Stdin)
				if err != nil {
					flog.Fatalf("read stdin: %v", err)
				}
				originalCommitMessage = string(byt)
			} else {
				var err error
				originalCommitMessage, err = readCommitFromCwd()
				if err != nil {
					flog.Fatalf("read commit: %v", err)
				}
			}

			originalCommitMessage = strings.TrimSpace(originalCommitMessage)

			// Remove existing emojis to make it easy to retry.
			originalCommitMessageNoEmoji := gomoji.RemoveEmojis(originalCommitMessage)

			resp, err := client.CreateCompletion(cmd.Context(), openai.CompletionRequest{
				Model:       "text-davinci-003",
				MaxTokens:   10,
				Temperature: 0.8,
				Prompt:      promptPrelude + "\nCommit: " + gomoji.RemoveEmojis(originalCommitMessageNoEmoji),
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
				color.Yellow("No emoji found, exiting\n")
				return
			}

			newCommitMessage := best + " " + originalCommitMessageNoEmoji

			if newCommitMessage == originalCommitMessage {
				color.Yellow("No change, exiting\n")
				return
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", newCommitMessage)
			if dryRun {
				color.Blue("Dry run, not amending\n")
			} else {
				if !adventureMode {
					p := promptui.Prompt{
						Label: "Commit? (y/n)",
					}
					result, err := p.Run()
					if err != nil {
						flog.Errorf("Prompt failed %v\n", err)
					}
					if result != "y" {
						color.Yellow("Exiting")
						os.Exit(1)
					}
				}

				color.Magenta("> git commit --amend \n")
				err = amendName(newCommitMessage)
				if err != nil {
					flog.Fatalf("amend: %v", err)
				}
			}
		},
	}

	cmd.Flags().BoolVar(&readStdin, "read-stdin", false, "read commit message from stdin")
	cmd.Flags().BoolVar(&dryRun, "dry", false, "don't actually amend the commit, just print the new message")
	cmd.Flags().BoolVarP(&adventureMode, "adventure", "y", false, "skip prompts, adventure mode")

	err := cmd.Execute()
	if err != nil {
		flog.Fatalf("%v", err)
	}
}
