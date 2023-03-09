.PHONY: gen

gen:
	go run ./cmd/gitemoji gen-fine-tuning > ./fine-tuning.json