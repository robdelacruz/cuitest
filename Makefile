all: t

dep:
	go get -u github.com/nsf/termbox-go

t: t.go widgets.go
	go build -o t t.go widgets.go

clean:
	rm -rf t

