SERVER_PORT=8080
TEMP_FILE=/tmp/metrics.json
#pwd for linux/apple
all: build vet iter1 iter2 iter3 iter4

build:
	go build -C ./cmd/agent main.go
	go build -C ./cmd/server main.go
vet:
	go vet -vettool=statictest ./...

iter1:
	./metricstest -test.v -test.run=^TestIteration1$$ \
                -binary-path=cmd/server/main.exe
iter2:
	./metricstest -test.v -test.run=^TestIteration2[AB]*$ \
                -source-path=. \
                -agent-binary-path=cmd/agent/main.exe
iter3:
	./metricstest -test.v -test.run=^TestIteration3[AB]*$ \
            -source-path=. \
            -binary-path=cmd/server/main.exe
            -agent-binary-path=cmd/agent/main.exe
iter4:
    SERVER_PORT=$(SERVER_PORT)
    ADDRESS="localhost:${SERVER_PORT}"
    TEMP_FILE=$(TEMP_FILE)
	./metricstest -test.v -test.run=^TestIteration4$ \
           -agent-binary-path=cmd/agent/main.exe \
           -binary-path=cmd/server/main.exe \
           -server-port=$SERVER_PORT \
           -source-path=.