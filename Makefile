SERVER_PORT=8080
#pwd for linux/apple
all: build vet iter1 iter2

build:
	cd
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