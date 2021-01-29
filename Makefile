build:
	go build -o packt main.go

run:
	go run main.go

install:
	go build -o packt main.go
	cp packt /usr/local/bin

clean:
	@rm packt