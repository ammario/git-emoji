# 🤪 git-emoji

Spruce up your commit messages with emojis courtesy of GPT-3 and [carloscuesta/gitmoji](https://github.com/carloscuesta/gitmoji/blob/master/packages/gitmojis/src/gitmojis.json).

## Usage

1. Sign up for [OpenAI](https://beta.openai.com/) and get an API key [here](https://platform.openai.com/account/api-keys).
1. Install the CLI with `go get github.com/ammario/git-emoji/cmd/gitemoji`.
1. Set your `OPENAI_API_KEY` environment variable to your API key.
1. Run `~/go/bin/gitemoji install`
1. The next time you're about to push a commit, run `git emoji` and follow the prompts.