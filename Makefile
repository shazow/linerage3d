BINARY = "linerage3d"

$(BINARY): *.go
	go build $(TAGS) .

run: $(BINARY)
	./$(BINARY)

test:
	go test

clean:
	rm $(BINARY)
