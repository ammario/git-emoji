package main

import (
	"encoding/json"
	"fmt"

	"github.com/coder/flog"
	"github.com/forPelevin/gomoji"
	"github.com/spf13/cobra"
)

type fineTuning struct {
	Prompt     string `json:"prompt,omitempty"`
	Completion string `json:"completion,omitempty"`
}

func genFineTuningsCmd() *cobra.Command {
	fineTunings := []fineTuning{
		{
			"Improve performance",
			`⚡️`,
		},
		{
			"Fix a bug",
			"🐛",
		},
		{
			"Fix a typo",
			"🔎",
		},
		{
			"Add a test",
			"🧪",
		},
		{
			"Add a feature",
			"🎁",
		},
	}
	return &cobra.Command{
		Use: "gen-fine-tunings",
		Run: func(cmd *cobra.Command, _ []string) {
			for _, ft := range fineTunings {
				emojis := gomoji.FindAll(ft.Completion)
				if len(emojis) == 0 {
					flog.Fatalf("no emojis found in %q", ft.Completion)
				} else if len(emojis) > 1 {
					flog.Fatalf("multiple emojis found in %q", ft.Completion)
				}
				byt, err := json.Marshal(ft)
				if err != nil {
					panic(err)
				}
				fmt.Printf("%s\n", byt)
			}
		},
	}
}
