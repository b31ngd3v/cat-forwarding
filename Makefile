BINARY = cat-forwarding
BIN_DIR = bin

build:
	@go build -o $(BIN_DIR)/$(BINARY)

clean:
	@rm -rf $(BIN_DIR)

test:
	@go test ./... -v
