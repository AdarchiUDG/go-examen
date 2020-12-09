package gui

import (
	"os"
	"github.com/jroimartin/gocui"
)

var gui *gocui.Gui

func GetGui() *gocui.Gui {
	return gui
}

func Layout(g *gocui.Gui) error {
	gui = g

	maxX, maxY := g.Size()
	if v, err := g.SetView("chat", 0, 0, maxX - 1, maxY - 4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Chat"
		v.Wrap = true
		v.Autoscroll = true
	}

	if v, err := g.SetView("input", 0, maxY - 3, maxX - 1, maxY - 1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Mensaje"
		v.Editable = true
		v.Editor = gocui.EditorFunc(editor)
		v.Highlight = true
		g.SetCurrentView("input")
	}
	return nil
}

func Quit(g *gocui.Gui, v *gocui.View) error {
	defer os.Exit(1)
	return gocui.ErrQuit
}

func editor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
		case ch != 0 && mod == 0:
			v.EditWrite(ch)
		case key == gocui.KeySpace:
			v.EditWrite(' ')
		case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
			v.EditDelete(true)
		case key == gocui.KeyDelete:
			v.EditDelete(false)
		case key == gocui.KeyInsert:
			v.Overwrite = !v.Overwrite
		case key == gocui.KeyArrowDown:
			v.MoveCursor(0, 1, false)
		case key == gocui.KeyArrowUp:
			v.MoveCursor(0, -1, false)
		case key == gocui.KeyArrowLeft:
			v.MoveCursor(-1, 0, false)
		case key == gocui.KeyArrowRight:
			v.MoveCursor(1, 0, false)
	}
}