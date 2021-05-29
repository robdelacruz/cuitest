SRCS = t.go tx.go txmenu.go txlistbox.go txlabel.go txtable.go txentry.go
all: t

dep:
	go get -u github.com/nsf/termbox-go

t: $(SRCS)
	go build -o t $(SRCS)

clean:
	rm -rf t

