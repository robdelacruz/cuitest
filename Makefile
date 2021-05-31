SRCS = t.go tx.go txmenu.go txlistbox.go txlabel.go txtable.go txentry.go
SRCS2 = db.go dbaccount.go dbcurrency.go
all: t

dep:
	go get -u github.com/nsf/termbox-go

t: $(SRCS) $(SRCS2)
	go build -o t $(SRCS) $(SRCS2)

clean:
	rm -rf t

