BINARY=linerage3d

$(BINARY): *.go
	go build $(TAGS) .

run: $(BINARY)
	./$(BINARY)

debug: clean
	$(eval TAGS=-tags gldebug)

test:
	go test

clean:
	rm -f $(BINARY)

.PHONY: debug
