bot.zip: bots/binary/data.go bots/binary/main.go
	cd bots; zip -9 -r ../bot.zip ./binary/

bots/process/main:
	# http://stackoverflow.com/a/21135705/907060
	env GOOS=linux GOARCH=amd64 go build -o bots/process/main -ldflags "-s" bots/process/main.go

bots/binary/data.go: bots/process/main
	go-bindata -nomemcopy -nocompress -o bots/binary/data.go  bots/process/main

bin/com:
	mkdir -p bin
	javac -d bin/ `find ./engine/blockbattle-engine/ -name '*.java'`
