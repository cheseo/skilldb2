.PHONY: run

run: skilldb2
	./skilldb2

skilldb2: *.go
	go build


