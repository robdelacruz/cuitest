all: t

dep:
	go get -u github.com/nsf/termbox-go

t: t.go widgets.go menuwidget.go listboxwidget.go
	go build -o t t.go widgets.go menuwidget.go listboxwidget.go

clean:
	rm -rf t

