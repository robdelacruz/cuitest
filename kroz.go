package main

import (
	"log"
	"os"

	tb "github.com/nsf/termbox-go"
)

var _log *log.Logger
var _termW, _termH int

func main() {
	flog, err := os.Create("./log.txt")
	if err != nil {
		panic(err)
	}
	defer flog.Close()
	_log = log.New(flog, "", 0)

	err = tb.Init()
	if err != nil {
		panic(err)
	}
	defer tb.Close()
	tb.SetOutputMode(tb.Output256)

	_termW, _termH = tb.Size()
	var x, y int
	x = 5
	y = 5
	var changed bool

	fg := tb.Attribute(156)
	bg := tb.Attribute(0)
	draw(x, y, fg, bg)

	chev := make(chan tb.Event)
	go pollEvents(chev)

	for {
		e := <-chev
		if e.Ch == 'q' {
			break
		}
		if e.Key == tb.KeyArrowUp {
			if y > 0 {
				y--
				changed = true
			}
		} else if e.Key == tb.KeyArrowDown {
			if y < _termH-1 {
				y++
				changed = true
			}
		} else if e.Key == tb.KeyArrowLeft {
			if x > 0 {
				x--
				changed = true
			}
		} else if e.Key == tb.KeyArrowRight {
			if x < _termW-1 {
				x++
				changed = true
			}
		}

		if changed {
			draw(x, y, fg, bg)
		}
	}
}

func draw(x, y int, fg, bg tb.Attribute) {
	tb.Clear(0, 0)
	tb.SetCell(x, y, '😀', fg, bg)
	tb.Flush()
}

func pollEvents(chev chan tb.Event) {
	for {
		e := tb.PollEvent()
		if e.Type != tb.EventKey {
			continue
		}
		chev <- e
	}
}
