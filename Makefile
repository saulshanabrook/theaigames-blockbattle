bot.zip: bots/binary/bot bots/binary/main.c
	cd bots; zip -9 -x "*.DS_Store" -x "binary/main"  -x "binary/a.out" -r ../bot.zip ./binary/

bots/binary/bot: $(shell find .  -iname "*.go" -type f) bots/process/nn.go
	# http://stackoverflow.com/a/21135705/907060
	env GOOS=linux GOARCH=amd64 go build -o bots/binary/bot bots/process/*.go

bots/process/nn.go: bots/process/nn
	go-bindata -nocompress -o bots/process/nn.go  bots/process/nn

bots/process/nn:
	go run rl/cmd/main.go

rl/engine/javac/com/theaigames/blockbattle/Blockbattle.class: $(shell find .  -iname "*.java" -type f)
	mkdir -p rl/engine/javac
	javac -d rl/engine/javac/ `find ./rl/engine/java -name '*.java'`

clean:
	rm -f bots/binary/bot bot.zip bots/process/nn.go

train: train/engine/javac/com/theaigames/blockbattle/Blockbattle.class
	go run train/engine/*.go

.PHONY: clean train
