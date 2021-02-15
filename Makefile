run-uncompiled:
	xvfb-run go run main.go 2>bot.log &
dependencies:
	go get github.com/360EntSecGroup-Skylar/excelize
	go get github.com/go-rod/rod
	go get github.com/go-telegram-bot-api/telegram-bot-api
build:
	go build -o sgobot main.go
run:
	xvfb-run ./sgobot 2>>bot.log &
config:
	read -p ""