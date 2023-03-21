lamport: lamport.go
		go build -o lamp.exe

clean:
		rm -rf lamport

all: lamport

.PHONY: all clean