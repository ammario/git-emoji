package main

import (
	_ "embed"
	"os"
	"time"

	"github.com/coder/flog"
	"github.com/coder/retry"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

//go:embed fine-tunings.jsonl
var fineTuningsRaw string

func newClient() *openai.Client {
	const keyEnv = "OPENAI_API_KEY"
	key := os.Getenv(keyEnv)
	if key == "" {
		flog.Fatalf("$%v is not set", keyEnv)
	}
	return openai.NewClient(key)
}

func installCmd() *cobra.Command {
	return &cobra.Command{
		Use: "install",
		Run: func(cmd *cobra.Command, _ []string) {
			client := newClient()
			tmpFi, err := os.CreateTemp("", "gitemoji")
			if err != nil {
				flog.Fatalf("create temp file: %v", err)
			}
			defer os.Remove(tmpFi.Name())

			err = os.WriteFile(tmpFi.Name(), []byte(fineTuningsRaw), 0o644)
			if err != nil {
				flog.Fatalf("write temp file: %v", err)
			}

			const fileName = "gitemoji"

			files, err := client.ListFiles(cmd.Context())
			if err != nil {
				flog.Fatalf("list files in openai: %v", err)
			}

			for _, file := range files.Files {
				if file.FileName == fileName {
					err = client.DeleteFile(cmd.Context(), file.ID)
					if err != nil {
						flog.Fatalf("delete file in openai: %v", err)
					}
					flog.Infof("old file %q deleted", file.ID)
				}
			}

			apiFile, err := client.CreateFile(cmd.Context(), openai.FileRequest{
				FileName: fileName,
				Purpose:  "fine-tune",
				FilePath: tmpFi.Name(),
			})
			if err != nil {
				flog.Fatalf("create file in openai: %v", err)
			}

			ftm, err := client.CreateFineTune(cmd.Context(), openai.FineTuneRequest{
				Model:        "davinci",
				TrainingFile: apiFile.ID,
			})
			if err != nil {
				flog.Fatalf("create fine-tune in openai: %v", err)
			}

			flog.Infof("fine-tune model %+v created", ftm)

			for r := retry.New(time.Second, time.Second*5); r.Wait(cmd.Context()); {
				ftm, err = client.GetFineTune(cmd.Context(), ftm.ID)
				if err != nil {
					flog.Fatalf("get fine-tune in openai: %v", err)
				}
				switch ftm.Status {
				case "succeeded":
					flog.Successf("fine-tune model %+v has succeeded", ftm)
					err = fineTuneModelID.Write([]byte(ftm.ID))
					if err != nil {
						flog.Fatalf("write fine-tune model id: %v", err)
					}
					return
				default:
					flog.Infof("fine-tune model %q is %q", ftm.ID, ftm.Status)
				}
			}
		},
	}
}
