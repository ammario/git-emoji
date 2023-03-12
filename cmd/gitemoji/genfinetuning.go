package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/forPelevin/gomoji"
)

type fineTuning struct {
	Prompt     string `json:"prompt,omitempty"`
	Completion string `json:"completion,omitempty"`
}

func generateFineTunings(w io.Writer) error {
	// Some inspiration taken from
	// https://github.com/carloscuesta/gitmoji/blob/master/packages/gitmojis/src/gitmojis.json.
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
		{
			"Add a dependency",
			"📦",
		},
		{
			"Remove a dependency",
			"🗑️",
		},
		{
			"Add documentation",
			"📝",
		},
		{
			"Refactor code",
			"♻️",
		},
		{
			"Improve structure / format of the code",
			"🎨",
		},
		{
			"Improve accessibility",
			"♿️",
		},
		{
			"Improve SEO",
			"📈",
		},
		{
			"Improve security",
			"🔒",
		},
		{
			"Fix major security issue",
			"🚨",
		},
		{
			"Critical hotfix",
			"🚑",
		},
		{
			"Update UI and style files",
			"🎨",
		},
		{
			"Remove dead code",
			"💀",
		},
		{
			"Add a logging message",
			"📟",
		},
	}

	for _, ft := range fineTunings {
		emojis := gomoji.FindAll(ft.Completion)
		if len(emojis) == 0 {
			return fmt.Errorf("no emojis found in %q", ft.Completion)
		} else if len(emojis) > 3 {
			return fmt.Errorf("too many emojis found in %q", ft.Completion)
		}
		ft.Prompt = promptPrelude + ft.Prompt
		byt, err := json.Marshal(ft)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "%s\n", byt)
		if err != nil {
			return err
		}
	}
	return nil
}
