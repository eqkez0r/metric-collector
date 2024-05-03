SERVER_PORT=8080
TEMP_FILE=/tmp/metrics.json
#pwd for linux/apple
all: build vet iter1 iter2 iter3 iter4 iter5

build:
	go build -C ./cmd/agent main.go
	go build -C ./cmd/server main.go
vet:
	go vet -vettool=statictest ./...

