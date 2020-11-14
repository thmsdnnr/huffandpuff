build:
	go build -o huff main.go

test:
	go test ./huffman
	./test.sh