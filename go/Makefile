all:
	$(MAKE) zverala
	$(MAKE) calendar/calendar

zverala: ./*.go
	go build -o ./bin/zverala

calendar/calendar: ./*.go
	go build -o ./bin/calendar ./calendar/