bot.zip: $(shell find bots/binary -type f) bots/binary/bot
	cd bots; zip -9 -x "*.DS_Store" -r ../bot.zip ./binary/

bots/binary/bot: $(shell find .  -iname "*.go" -type f)
	env GOOS=linux GOARCH=amd64 go build -gcflags=-l  -ldflags "-s -extldflags \"-static\"" -o bots/binary/bot bots/process/*.go

rl/engine/javac/com/theaigames/blockbattle/Blockbattle.class: $(shell find .  -iname "*.java" -type f)
	mkdir -p rl/engine/javac
	javac -d rl/engine/javac/ `find ./rl/engine/java -name '*.java'`

clean:
	rm -f bots/binary/bot bot.zip bots/process/nn.go

online: train/engine/javac/com/theaigames/blockbattle/Blockbattle.class
	go run rl/online/*.go

.PHONY: clean train
