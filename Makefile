all: fmt vet lint

fmt:
	go fmt .

vet:
	go vet .

lint:
	golint .

install:
	go install ./god
