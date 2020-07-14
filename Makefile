.PHONY: build
build:
		go build -v .

.PHONY: run
run:
		go build -v .
		delau-notify-bot

.DEFAULT_GOAL := build 