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
	window.SetDefaultSize(480, 320)              //设置大小

	//鼠标按下事件处理
	window.Connect("button-press-event", func(ctx *glib.CallbackContext) {
		//获取鼠键按下属性结构体变量，系统内部的变量，不是用户传参变量
		arg := ctx.Args(0)
		event := *(**gdk.EventButton)(unsafe.Pointer(&arg))

		if event.Type == int(gdk.BUTTON_PRESS) { //单击
			fmt.Printf("单击, ")
		} else if event.Type == int(gdk.BUTTON2_PRESS) { //双击
			fmt.Printf("双击, ")
		}

		//fmt.Println("button = ", event.Button)
		if event.Button == 1 {
			fmt.Printf("左键, ")
		} else if event.Button == 2 {
			fmt.Printf("中间键, ")
		} else if event.Button == 3 {
			fmt.Printf("右键, ")
		}

		fmt.Println("坐标:", int(event.X), int(event.Y))
	})

	//鼠标移动事件处理
	window.Connect("motion-notify-event", func(ctx *glib.CallbackContext) {
		//获取鼠标移动属性结构体变量，系统内部的变量，不是用户传参变量
		arg := ctx.Args(0)
		event := *(**gdk.EventMotion)(unsafe.Pointer(&arg))

		fmt.Println("移动坐标:", int(event.X), int(event.Y))
	})

	//添加鼠标按下事件
	//BUTTON_PRESS_MASK: 鼠标按下，触发信号"button-press-event"
	//BUTTON_RELEASE_MASK：鼠标抬起，触发"button-release-event"
	//鼠标移动都是触发"motion-notify-event"
	//BUTTON_MOTION_MASK: 鼠标移动，按下任何键移动都可以
	//BUTTON1_MOTION_MASK：鼠标移动，按住左键移动才触发
	//BUTTON2_MOTION_MASK：鼠标移动，按住中间键移动才触发
	//BUTTON3_MOTION_MASK：鼠标移动，按住右键移动才触发
	window.SetEvents(int(gdk.BUTTON_PRESS_MASK | gdk.BUTTON1_MOTION_MASK))

	window.Connect("destroy", gtk.MainQuit) //关闭窗口
	window.ShowAll()                        //显示控件

	gtk.Main()
}

//标事件，可以理解为操作鼠标的动作。对于窗口而言，用户操作鼠标，窗口检测到鼠标的操作( 产生一个信号 )，然后去做相应处理( 调用其规定的回调函数 )，即为鼠标事件。

//窗口默认不捕获鼠标的操作，需要手动添加让其捕获：

//    //添加鼠标按下事件
//    //BUTTON_PRESS_MASK: 鼠标按下，触发信号"button-press-event"
//    //BUTTON_RELEASE_MASK：鼠标抬起，触发"button-release-event"
//    //鼠标移动都是触发"motion-notify-event"
//    //BUTTON_MOTION_MASK: 鼠标移动，按下任何键移动都可以
//    //BUTTON1_MOTION_MASK：鼠标移动，按住左键移动才触发
//    //BUTTON2_MOTION_MASK：鼠标移动，按住中间键移动才触发
//    //BUTTON3_MOTION_MASK：鼠标移动，按住右键移动才触发
//    window.SetEvents(int(gdk.BUTTON_PRESS_MASK | gdk.BUTTON1_MOTION_MASK))

//事件标识 	触发条件
//“button-press-event” 	鼠标按下时触发
//“button-release-event” 	鼠标抬起时触发
//“motion-notify-event” 	按住鼠标移动时触发
