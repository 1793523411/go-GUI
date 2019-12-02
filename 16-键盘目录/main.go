package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

func main() {
	gtk.Init(&os.Args)

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL) //新建窗口
	window.SetPosition(gtk.WIN_POS_CENTER)       //默认居中显示
	window.SetDefaultSize(480, 320)              //设置大小

	//键盘按下事件处理
	window.Connect("key-press-event", func(ctx *glib.CallbackContext) {
		//获取键盘按下属性结构体变量，系统内部的变量，不是用户传参变量
		arg := ctx.Args(0)
		event := *(**gdk.EventKey)(unsafe.Pointer(&arg))

		//event.Keyval：获取按下(释放)键盘键值，每个键值对于一个ASCII码
		key := event.Keyval
		if gdk.KEY_Up == key {
			fmt.Println("上")
		} else if gdk.KEY_Down == key {
			fmt.Println("下")
		} else if gdk.KEY_Left == key {
			fmt.Println("左")
		} else if gdk.KEY_Right == key {
			fmt.Println("右")
		}

		fmt.Println("key = ", event.Keyval)

	})

	window.Connect("destroy", gtk.MainQuit) //关闭窗口

	window.ShowAll() //显示控件

	gtk.Main()
}

//事件标识 	触发条件
//“key-press-event” 	键盘按下时触发
//“key-release-event” 	键盘抬起时触发
