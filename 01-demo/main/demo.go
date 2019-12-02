package main

import (
	"fmt"
	"os"

	"github.com/mattn/go-gtk/gtk"
)

func main() {
	gtk.Init(&os.Args)

	win := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	win.SetTitle("go gtk")
	win.SetSizeRequest(480, 320)
	win.Show()

	fmt.Println("before")
	gtk.Main()
	fmt.Println("over")
}
