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
	builder.AddFromFile("test.glade") //读取glade文件

	// 获取窗口控件指针，注意"window1"要和glade里的标志名称匹配
	window := gtk.WindowFromObject(builder.GetObject("window1"))

	//获取image控件
	image1 := gtk.ImageFromObject(builder.GetObject("image1"))

	//获取image控件大小
	var w, h int
	image1.GetSizeRequest(&w, &h)
	fmt.Println(w, h)

	//创建pixbuf，指定大小（宽度和高度），image有多大就设置多大
	//最后一个参数false代表不保存图片原来的尺寸
	pixbuf1, _ := gdkpixbuf.NewPixbufFromFileAtScale("image/face.png", w, h, false)

	//image设置pixbuf
	image1.SetFromPixbuf(pixbuf1)

	//pixbuf1使用完毕，需要释放资源
	pixbuf1.Unref()

	//按窗口关闭按钮，自动触发"destroy"信号
	window.Connect("destroy", gtk.MainQuit)

	window.Show()

	gtk.Main()
}
