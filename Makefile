BINARY = "gogame"

$(BINARY): *.go
	go build .

run:
	go run *.go

clean:
	rm $(BINARY)
