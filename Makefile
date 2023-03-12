.PHONY: gen

gen:
	go run ./cmd/gitemoji gen-fine-tunings > ./cmd/fine-tunings.jsonl