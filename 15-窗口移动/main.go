package main

import (
	//  "fmt"
	"os"
	"unsafe"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

func main() {
	gtk.Init(&os.Args)

	window := gtk.NewWindow(gtk.WINDOW_POPUP) //新建窗口，弹出类型，即为无边框窗口
	window.SetPosition(gtk.WIN_POS_CENTER)    //默认居中显示
	window.SetDefaultSize(480, 320)           //设置大小

	var x, y int

	//鼠标按下事件处理
	window.Connect("button-press-event", func(ctx *glib.CallbackContext) {
		//获取鼠键按下属性结构体变量，系统内部的变量，不是用户传参变量
		arg := ctx.Args(0)
		event := *(**gdk.EventButton)(unsafe.Pointer(&arg))

		if event.Button == 1 { //左键
			x, y = int(event.X), int(event.Y) //保存点击的起点坐标
		} else if event.Button == 3 { //右键
			//右键，关闭窗口
			gtk.MainQuit()
		}
	})

	//鼠标移动事件处理
	window.Connect("motion-notify-event", func(ctx *glib.CallbackContext) {
		//获取鼠标属性结构体变量，系统内部的变量，不是用户传参变量
		arg := ctx.Args(0)
		//还是EventButton
		event := *(**gdk.EventButton)(unsafe.Pointer(&arg))

		//移动窗口
		window.Move(int(event.XRoot)-x, int(event.YRoot)-y)
	})

	//添加鼠标按下事件
	window.SetEvents(int(gdk.BUTTON_PRESS_MASK | gdk.BUTTON1_MOTION_MASK))

	window.ShowAll() //显示控件

	gtk.Main()
}
