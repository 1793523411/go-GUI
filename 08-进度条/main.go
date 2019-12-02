package main

import (
	"fmt"
	"os"

	"github.com/mattn/go-gtk/gtk"
)

func main() {
	gtk.Init(&os.Args)

	builder := gtk.NewBuilder()       //新建builder
	builder.AddFromFile("test.glade") //读取glade文件

	// 获取窗口控件指针，注意"window1"要和glade里的标志名称匹配
	window := gtk.WindowFromObject(builder.GetObject("window1"))
	window.SetSizeRequest(440, 250) //设置窗口大小

	//获取进度条控件
	pb := gtk.ProgressBarFromObject(builder.GetObject("progressbar1"))

	//设置进度, 范围0.0~1.0
	pb.SetFraction(0.5)

	//设置进度条显示的文本
	pb.SetText("50%")

	//获取进度
	value := pb.GetFraction()
	fmt.Println("value = ", value)

	//按窗口关闭按钮，自动触发"destroy"信号
	window.Connect("destroy", gtk.MainQuit)

	window.Show()

	gtk.Main()
}
