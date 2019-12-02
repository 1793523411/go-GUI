package main

import (
	"os"

	"github.com/mattn/go-gtk/gtk"
)

func main() {
	gtk.Init(&os.Args)

	builder := gtk.NewBuilder()       //新建builder
	builder.AddFromFile("test.glade") //读取glade文件

	// 获取窗口控件指针，注意"window1"要和glade里的标志名称匹配
	window := gtk.WindowFromObject(builder.GetObject("window1"))

	window.SetSizeRequest(300, 240)        //设置窗口大小
	window.SetTitle("hello go")            //设置标题
	window.SetIconFromFile("face.png")     //设置icon
	window.SetResizable(false)             //设置不可伸缩
	window.SetPosition(gtk.WIN_POS_CENTER) //设置居中显示

	//按窗口关闭按钮，自动触发"destroy"信号
	window.Connect("destroy", gtk.MainQuit)

	window.ShowAll()

	gtk.Main()
}

//信号标识	触发条件
//“destroy”	按关闭窗口按钮时触发
