BINARY = "gogame"

$(BINARY): *.go
	go build .

clean:
	rm $(BINARY)
