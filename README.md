# 🤪 git-emoji

Spruce up your commit messages with emojis courtesy of GPT-3.

## Usage

1. Sign up for [OpenAI](https://beta.openai.com/) and get an API key [here](https://platform.openai.com/account/api-keys).
1. Install the CLI with `go install github.com/ammario/git-emoji/cmd/gitemoji`.
1. Set your `OPENAI_API_KEY` environment variable to your API key.
1. The next time you're about to push a commit, run `gitemoji` right before and follow the prompts.

## Fully Commit

If you care deeply about emojis, you can add the following to your `~/.gitconfig`:

```
[alias]
    emoji = "! go run github.com/ammario/git-emoji/cmd/gitemoji"
```