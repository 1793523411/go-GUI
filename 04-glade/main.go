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
	b1 := gtk.ButtonFromObject(builder.GetObject("123456"))        //获取按钮1
	b2 := gtk.ButtonFromObject(builder.GetObject("togglebutton1")) //获取按钮2

	//信号处理
	b1.Connect("clicked", func() {
		//获取按钮内容
		fmt.Println("button txt = ", b1.GetLabel())
	})

	b2.Connect("clicked", func() {
		//获取按钮内容
		fmt.Println("button txt = ", b2.GetLabel())
		gtk.MainQuit() //关闭窗口
	})

	//按窗口关闭按钮，自动触发"destroy"信号
	window.Connect("destroy", gtk.MainQuit)

	window.ShowAll()

	gtk.Main()
}

//可以简单分为两步：

//1）读取glade文件

//// 创建GtkBuilder对象,GtkBuilder在<gtk/gtk.h>声明

//GtkBuilder *builder = gtk_builder_new();

//// 读取test.glade文件的信息，保存在builder指针变量里

//gtk_builder_add_from_file(builder, "./test.glade", NULL);

//2）获取glade文件里的控件

//// 获取窗口控件指针，注意"window1" 要和glade里的标志名称匹配

//GtkWidget *window = GTK_WIDGET(gtk_builder_get_object (builder, "window1"));
