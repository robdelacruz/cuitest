SRCS = t.go waccounts.go
SRCS2 = tx.go txmenu.go txlistbox.go txlabel.go txtable.go txentry.go
SRCS3 = db.go dbaccount.go dbcurrency.go
all: t

dep:
	go get -u github.com/nsf/termbox-go

t: $(SRCS) $(SRCS2) $(SRCS3)
	go build -o t $(SRCS) $(SRCS2) $(SRCS3)

clean:
	rm -rf t

