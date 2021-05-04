all: t

dep:
	go get -u github.com/nsf/termbox-go

t: t.go
	go build -o t t.go

clean:
	rm -rf t

