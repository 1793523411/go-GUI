package main

import (
	"fmt"
	"os"

	"github.com/mattn/go-gtk/gtk"
)

func main() {
	gtk.Init(&os.Args)

	builder := gtk.NewBuilder()          //新建builder
	builder.AddFromFile("../test.glade") //读取glade文件

	// 获取窗口控件指针，注意"window1"要和glade里的标志名称匹配
	window := gtk.WindowFromObject(builder.GetObject("window1"))

	//获取label控件
	labelOne := gtk.LabelFromObject(builder.GetObject("label1"))
	labelTwo := gtk.LabelFromObject(builder.GetObject("label2"))

	fmt.Println("labelOne = ", labelOne.GetText()) //获取label内容
	labelOne.SetText("你大爷")                        //设置内容
	labelTwo.SetText("^_^")

	//按窗口关闭按钮，自动触发"destroy"信号
	window.Connect("destroy", gtk.MainQuit)

	window.Show()

	gtk.Main()
}
