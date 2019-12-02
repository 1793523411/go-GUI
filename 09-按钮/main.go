package main

import (
	"fmt"
	"os"

	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/gtk"
)

func main() {
	gtk.Init(&os.Args)

	builder := gtk.NewBuilder()       //新建builder
	builder.AddFromFile("tset.glade") //读取glade文件

	// 获取窗口控件指针，注意"window1"要和glade里的标志名称匹配
	window := gtk.WindowFromObject(builder.GetObject("window1"))

	//获取button控件
	button1 := gtk.ButtonFromObject(builder.GetObject("button1"))
	button2 := gtk.ButtonFromObject(builder.GetObject("button2"))

	button1.SetLabel("@_~")                           //按钮设置文本信息
	fmt.Println("button1 txt = ", button1.GetLabel()) //获取按钮内容
	button1.SetSensitive(false)                       //按钮变灰色，不能按

	//获取按钮2的大小
	var w, h int
	button2.GetSizeRequest(&w, &h)
	fmt.Println(w, h)

	//创建pixbuf，指定大小（宽度和高度
	pixbuf, _ := gdkpixbuf.NewPixbufFromFileAtScale("head.png", w-10, h-10, false)

	//通过pixbuf新建image
	image := gtk.NewImageFromPixbuf(pixbuf)

	//释放pixbuf资源
	pixbuf.Unref()

	//按钮设置image
	button2.SetImage(image)

	//按钮信号处理
	button2.Connect("clicked", func() {
		fmt.Println("按钮2被按下")
	})

	//按窗口关闭按钮，自动触发"destroy"信号
	window.Connect("destroy", gtk.MainQuit)

	window.ShowAll()

	gtk.Main()
}

//信号标识 	触发条件
//“clicked” 	按下按钮时触发
//“pressed” 	按下按钮时触发
//“released” 	释放按钮时触发
