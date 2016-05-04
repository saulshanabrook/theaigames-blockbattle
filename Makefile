bot.zip: bots/binary/bot bots/binary/main.c
	cd bots; zip -9 -x "*.DS_Store" -x "binary/main" -r ../bot.zip ./binary/

bots/binary/bot: $(shell find .  -iname "*.go" -type f)
	# http://stackoverflow.com/a/21135705/907060
	env GOOS=linux GOARCH=amd64 go build -o bots/binary/bot -ldflags "-s" bots/process/*.go

rl/engine/javac/com/theaigames/blockbattle/Blockbattle.class: $(shell find .  -iname "*.java" -type f)
	mkdir -p rl/engine/javac
	javac -d rl/engine/javac/ `find ./rl/engine/java -name '*.java'`

clean:
	rm -f bots/binary/bot bot.zip

train: train/engine/javac/com/theaigames/blockbattle/Blockbattle.class
	go run train/engine/*.go

.PHONY: clean train
