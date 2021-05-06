all: t

dep:
	go get -u github.com/nsf/termbox-go

t: t.go widgets.go menuwidget.go
	go build -o t t.go widgets.go menuwidget.go

clean:
	rm -rf t

