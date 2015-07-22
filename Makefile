BINARY = "linerage3d"

$(BINARY): *.go
	go build .

run: $(BINARY)
	./$(BINARY)

test:
	go test

clean:
	rm $(BINARY)
