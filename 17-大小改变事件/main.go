package main

import (
	"fmt"
	"os"

	"github.com/mattn/go-gtk/gtk"
)

func main() {
	gtk.Init(&os.Args)

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL) //新建窗口
	window.SetPosition(gtk.WIN_POS_CENTER)       //默认居中显示
	window.SetDefaultSize(480, 320)              //设置大小

	//大小改变事件，改变窗口大小，自动触发此事件
	window.Connect("configure_event", func() {
		//获取窗口尺寸
		var w, h int
		window.GetSize(&w, &h)
		fmt.Println(w, h)
	})

	window.Connect("destroy", gtk.MainQuit) //关闭窗口

	window.ShowAll() //显示控件

	gtk.Main()
}
